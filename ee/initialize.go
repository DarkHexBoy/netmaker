//go:build ee
// +build ee

package ee

import (
	controller "github.com/gravitl/netmaker/controllers"
	"github.com/gravitl/netmaker/ee/ee_controllers"
	eelogic "github.com/gravitl/netmaker/ee/logic"
	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/logic"
	"github.com/gravitl/netmaker/models"
	"github.com/gravitl/netmaker/servercfg"
)

// InitEE - Initialize EE Logic
func InitEE() {
	setIsEnterprise()
	models.SetLogo(retrieveEELogo())
	controller.HttpHandlers = append(controller.HttpHandlers, ee_controllers.MetricHandlers)
	logic.EnterpriseCheckFuncs = append(logic.EnterpriseCheckFuncs, func() {
		// == License Handling ==
		ValidateLicense()
		if Limits.FreeTier {
			logger.Log(0, "proceeding with Free Tier license")
		} else {
			logger.Log(0, "proceeding with Paid Tier license")
		}
		// == End License Handling ==
		AddLicenseHooks()
	})
	logic.EnterpriseFailoverFunc = eelogic.SetFailover
	logic.EnterpriseResetFailoverFunc = eelogic.ResetFailover
	resetFailover()
}

func setControllerLimits() {
	logic.Node_Limit = Limits.Nodes
	logic.Users_Limit = Limits.Users
	logic.Clients_Limit = Limits.Clients
	logic.Free_Tier = Limits.FreeTier
	servercfg.Is_EE = true
}

func resetFailover() {
	nets, err := logic.GetNetworks()
	if err == nil {
		for _, net := range nets {
			err = logic.EnterpriseResetFailoverFunc.(func(string) error)(net.NetID)
			if err != nil {
				logger.Log(0, "failed to reset failover on network", net.NetID, ":", err.Error())
			}
		}
	}
}

func retrieveEELogo() string {
	return `              
 __   __     ______     ______   __    __     ______     __  __     ______     ______    
/\ "-.\ \   /\  ___\   /\__  _\ /\ "-./  \   /\  __ \   /\ \/ /    /\  ___\   /\  == \   
\ \ \-.  \  \ \  __\   \/_/\ \/ \ \ \-./\ \  \ \  __ \  \ \  _"-.  \ \  __\   \ \  __<   
 \ \_\\"\_\  \ \_____\    \ \_\  \ \_\ \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_\ \_\ 
  \/_/ \/_/   \/_____/     \/_/   \/_/  \/_/   \/_/\/_/   \/_/\/_/   \/_____/   \/_/ /_/ 
                                                                                         																							 
                                   ___    ___   ____                        
           ____  ____  ____       / _ \  / _ \ / __ \       ____  ____  ____
          /___/ /___/ /___/      / ___/ / , _// /_/ /      /___/ /___/ /___/
         /___/ /___/ /___/      /_/    /_/|_| \____/      /___/ /___/ /___/ 
                                                                            
`
}
