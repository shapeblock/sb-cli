package cmd

import (
	//"fmt"
	//"net/http"
	//"os"

	//"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

type VolumeDeletePayload struct {
	Volumes []string `json:"delete"`
}
func fetchVolumeValue(volVars []*VolumeSelect) [] string {
	var vols []string
	for _,volume:=range volVars{
		vols=append(vols, volume.Name)
	}
	return vols
}

var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an volume",
	Run:   volumeDelete,
}

func volumeDelete(cmd *cobra.Command, args []string) {
//	apps, err := fetchApps()
//	if err != nil {
		//fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
//	}

	//app := selectApp(apps)

	//appDetail, err := fetchAppDetail(app.UUID)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}




func init() {
	appVolumeCmd.AddCommand(volumeDeleteCmd)
}
