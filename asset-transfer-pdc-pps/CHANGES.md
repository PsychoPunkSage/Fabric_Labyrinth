1. Create a project `vehicle lifecycle management` with two org
2. Use PDC to share data
3. create pdc data in org1 and access in org2
4. Implement chaincode lifecycle and pdc  


Vehicle Lifetime Management::>
+ Type::>     Vehicle Name        === Name   
+ ID::>       Vehicle number      === RegNum
+ Color::>    Vehicle Company     === Company
+ Size::>     Vehicle Year of Reg === MfgYear
+ Owner::>    Vehicle Owner       === Owner

+ ID::>       Vehicle number      === RegNum
+ AppraisedValue::> Vehicle Life  === Life











type Asset struct {
	+ Type  string `json:"objectType"` 
	+ ID    string `json:"assetID"`
	+ Color string `json:"color"`
    + Size  int    `json:"size"`
	+ Owner string `json:"owner"`  
}

// AssetPrivateDetails describes details that are private to owners
type AssetPrivateDetails struct {
	ID             string `json:"assetID"`
	AppraisedValue int    `json:"appraisedValue"`
}

// TransferAgreement describes the buyer agreement returned by ReadTransferAgreement
type TransferAgreement struct {
	ID      string `json:"assetID"`
	BuyerID string `json:"buyerID"`
}