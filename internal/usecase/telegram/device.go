package telegram

import (
	"fmt"
	"miner-fetch/internal/device"
	"miner-fetch/internal/storage"
)

type Usecase struct {
	Storage *storage.Storage
}

func (u *Usecase) GetDevicesInfo() string {
	devices := u.Storage.GetDevices()

	var text string

	for _, d := range devices {
		versionCommand := device.VersionCommand{}
		err := d.SendCommand(&versionCommand)
		if err != nil {
			fmt.Println(err)
		}

		statsCommand := device.StatsCommand{}
		err = d.SendCommand(&statsCommand)
		if err != nil {
			fmt.Println(err)
		}

		poolsCommand := device.PoolsCommand{}
		err = d.SendCommand(&poolsCommand)
		if err != nil {
			fmt.Println(err)
		}

		text += fmt.Sprintf(
			"%s [%s]\n",
			versionCommand.Response.Version[0].Type,
			poolsCommand.Response.Pools[0].User,
		)

		text += fmt.Sprintf(
			"Temp 1 — %d %d\n",
			statsCommand.Response.Stats[1].Temp1,
			statsCommand.Response.Stats[1].Temp21,
		)

		text += fmt.Sprintf(
			"Temp 2 — %d %d\n",
			statsCommand.Response.Stats[1].Temp2,
			statsCommand.Response.Stats[1].Temp22,
		)

		text += fmt.Sprintf(
			"Temp 3 — %d %d\n",
			statsCommand.Response.Stats[1].Temp3,
			statsCommand.Response.Stats[1].Temp23,
		)

		text += "\n"
	}

	return text
}
