package node

import (
	"strconv"

	"ad2nx/api/panel"
	"ad2nx/common/sysinfo"

	log "github.com/sirupsen/logrus"
)

func (c *Controller) reportUserTrafficTask() (err error) {
	req := &panel.NodeReportRequest{}

	// 1. 采集 traffic
	userTraffic, _ := c.server.GetUserTrafficSlice(c.tag, true)
	var totalUp, totalDown int64
	if len(userTraffic) > 0 {
		req.Traffic = make(map[int][]int64, len(userTraffic))
		for _, t := range userTraffic {
			req.Traffic[t.UID] = []int64{t.Upload, t.Download}
			totalUp += t.Upload
			totalDown += t.Download
		}
	}

	// 2. 采集 alive (在线 IP) 和 online (连接数)
	if onlineDevice, err := c.limiter.GetOnlineDevice(); err == nil && len(*onlineDevice) > 0 {
		// 过滤低流量用户（保留原有逻辑）
		nocountUID := make(map[int]struct{})
		for _, traffic := range userTraffic {
			total := traffic.Upload + traffic.Download
			if total < int64(c.Options.DeviceOnlineMinTraffic*1000) {
				nocountUID[traffic.UID] = struct{}{}
			}
		}

		alive := make(map[int][]string)
		online := make(map[int]int)
		for _, ou := range *onlineDevice {
			// online 统计所有连接数（不过滤）
			online[ou.UID]++
			// alive 过滤低流量用户
			if _, skip := nocountUID[ou.UID]; !skip {
				alive[ou.UID] = append(alive[ou.UID], ou.IP)
			}
		}
		if len(alive) > 0 {
			req.Alive = alive
		}
		if len(online) > 0 {
			req.Online = online
		}

		log.WithField("tag", c.tag).Infof("Total %d online users, %d reported in alive",
			len(*onlineDevice), len(alive))
	}

	// 3. 采集 status（CPU/Mem/Swap/Disk + 带宽）
	interval := c.info.PushInterval.Seconds()
	if interval <= 0 {
		interval = 60 // fallback
	}
	status := sysinfo.GetNodeStatus()
	status.InboundSpeed = int64(float64(totalUp) / interval)
	status.OutboundSpeed = int64(float64(totalDown) / interval)
	req.Status = status

	// 4. 采集 metrics
	activeUsers := len(userTraffic)
	req.Metrics = sysinfo.GetNodeMetrics(
		len(c.userList),
		activeUsers,
		status.InboundSpeed,
		status.OutboundSpeed,
	)

	// 统一上报
	if err = c.apiClient.ReportNodeData(req); err != nil {
		log.WithFields(log.Fields{
			"tag": c.tag,
			"err": err,
		}).Error("Report node data failed")
	} else {
		log.WithField("tag", c.tag).Infof("Report %d users traffic, speed: ↑%d ↓%d bytes/s",
			len(userTraffic), status.InboundSpeed, status.OutboundSpeed)
	}

	userTraffic = nil
	return nil
}

func compareUserList(old, new []panel.UserInfo) (deleted, added []panel.UserInfo) {
	oldMap := make(map[string]int)
	for i, user := range old {
		key := user.Uuid + strconv.Itoa(user.SpeedLimit)
		oldMap[key] = i
	}

	for _, user := range new {
		key := user.Uuid + strconv.Itoa(user.SpeedLimit)
		if _, exists := oldMap[key]; !exists {
			added = append(added, user)
		} else {
			delete(oldMap, key)
		}
	}

	for _, index := range oldMap {
		deleted = append(deleted, old[index])
	}

	return deleted, added
}
