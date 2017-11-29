package types

import (
	"time"
)

// MMEConf parameters for the ENB to connect to
type MMEConf	struct{
	Ipv4		string	`json:"ipv4"`
	Ipv6		string	`json:"ipv6"`
	Preference	string	`json:"preference"`
}

//EnbIF is the ENB_INTERFACE section
type EnbIF struct{
	S1MmeIF		string		`json:"S1_MME_IF,omitempty"`
	S1MmeAddr	string		`json:"S1_MME_Addr,omitempty"`
	S1UIF		string		`json:"S1_U_IF,omitempty"`
	S1UAddr		string		`json:"S1_U_Addr,omitempty"`
	S1UPort		int			`json:"S1_U_Port,omitempty"`
}

//OAIEnb contains typical OAI eNodeB parameters
type OAIEnb struct{
	ID		string	`json:"enb_id"`
	Cell	string	`json:"cell_type"`
	Name	string	`json:"enb_name"`
	Tra		int		`json:"enb_tra"`
	Mcc		int		`json:"enb_mcc"`
	Mnc		int		`json:"enb_mnc"`
	Dl		int		`json:"enb_dl"`
	Ul		int		`json:"enb_ul"`
	Tx		int		`json:"nb_tx"`
	Rx		int		`json:"nb_rx"`
	Mme		MMEConf	`json:"mme"`
	IF		EnbIF	`json:"-"`
}
/*	In the case you wanted other sections for your configuration,
	you should define your parameters in a new struct and include
	it into the VNFParams struct like I did for OAIEnb. This way
	you could extend your configuration by extending this struct 
	alone*/

//VNFParams contains the parameters passed to the VNF. Could be extended at will.
type VNFParams struct{
	Enb		OAIEnb		`json:"enb"`
}

//Config is for saving in the storage and generating the actual config for the VNF
type Config struct {
	ID			string					`json:"id,omitempty"`
	Alias		string					`json:"alias"`
	Session		string					`json:"session_id"`
	Version		string					`json:"template_version"`
	//Instances	map[string]*Instance	`json:"instance_id"`
	CreatedAt	time.Time   			`json:"created_at"`
	Tags		[]string				`json:"tags"`
	Content		*VNFParams				`json:"params"`
}