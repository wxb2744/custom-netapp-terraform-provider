package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

type apiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// only list what is needed
type workingEnvironmentInfo struct {
	Name                   string `json:"name"`
	PublicID               string `json:"publicId"`
	CloudProviderName      string `json:"cloudProviderName"`
	IsHA                   bool   `json:"isHA"`
	WorkingEnvironmentType string `json:"workingEnvironmentType"`
	SvmName                string `json:"svmName"`
}

type workingEnvironmentResult struct {
	VsaWorkingEnvironment       []workingEnvironmentInfo `json:"vsaWorkingEnvironments"`
	OnPremWorkingEnvironments   []workingEnvironmentInfo `json:"onPremWorkingEnvironments"`
	AzureVsaWorkingEnvironments []workingEnvironmentInfo `json:"azureVsaWorkingEnvironments"`
	GcpVsaWorkingEnvironments   []workingEnvironmentInfo `json:"gcpVsaWorkingEnvironments"`
}

type workingEnvironmentOntapClusterPropertiesResponse struct {
	ActionsRequired                interface{}            `json:"actionsRequired"`
	ActiveActions                  interface{}            `json:"activeActions"`
	AwsProperties                  interface{}            `json:"awsProperties"` // aws
	CapacityFeatures               interface{}            `json:"capacityFeatures"`
	CbsProperties                  interface{}            `json:"cbsProperties"`
	CloudSyncProperties            interface{}            `json:"cloudSyncProperties"` // aws
	CloudProviderName              string                 `json:"cloudProviderName"`
	ComplianceProperties           interface{}            `json:"complianceProperties"`
	CreatorUserEmail               string                 `json:"creatorUserEmail"`
	CronJobSchedules               interface{}            `json:"cronJobSchedules"` // aws
	EncryptionProperties           interface{}            `json:"encryptionProperties"`
	FpolicyProperties              interface{}            `json:"fpolicyProperties"`
	HAProperties                   interface{}            `json:"haProperties"`
	InterClusterLifs               interface{}            `json:"interClusterLifs"` // aws
	IsHA                           bool                   `json:"isHA"`
	LicensesInformation            interface{}            `json:"licensesInformation"`
	MonitoringProperties           interface{}            `json:"monitoringProperties"`
	Name                           string                 `json:"name"`
	OntapClusterProperties         ontapClusterProperties `json:"ontapClusterProperties"`
	ProviderProperties             interface{}            `json:"providerProperties"`
	PublicID                       string                 `json:"publicId"`
	ReplicationProperties          interface{}            `json:"replicationProperties"`
	ReservedSize                   interface{}            `json:"reservedSize"`
	SaasProperties                 interface{}            `json:"saasProperties"`
	Schedules                      interface{}            `json:"schedules"`
	SnapshotPolicies               interface{}            `json:"snapshotPolicies"`
	Status                         cvoStatus              `json:"status"`
	SupportRegistrationInformation []interface{}          `json:"supportRegistrationInformation"`
	SupportRegistrationProperties  interface{}            `json:"supportRegistrationProperties"`
	SupportedFeatures              interface{}            `json:"supportedFeatures"`
	SvmName                        string                 `json:"svmName"`
	Svms                           interface{}            `json:"svms"`
	TenantID                       string                 `json:"tenantId"`
	WorkingEnvironmentType         string                 `json:"workingEnvironmentType"`
}

type ontapClusterProperties struct {
	BroadcastDomainInfo              []broadcastDomainInfo `json:"broadcastDomainInfo"`
	CanConfigureCapacityTier         bool                  `json:"canConfigureCapacityTier"`
	CapacityTierInfo                 capacityTierInfo      `json:"capacityTierInfo"`
	ClusterName                      string                `json:"clusterName"`
	ClusterUUID                      string                `json:"clusterUuid"`
	CreationTime                     interface{}           `json:"creationTime"`
	Evaluation                       bool                  `json:"evaluation"`
	IsSpaceReportingLogical          bool                  `json:"isSpaceReportingLogical"`
	LastModifiedOffbox               interface{}           `json:"lastModifiedOffbox"`
	LicensePackageName               interface{}           `json:"licensePackageName"`
	LicenseType                      licenseType           `json:"licenseType"`
	Nodes                            []node                `json:"nodes"`
	OffboxTarget                     bool                  `json:"offboxTarget"`
	OntapVersion                     string                `json:"ontapVersion"`
	SystemManagerURL                 string                `json:"systemManagerUrl"`
	UpgradeVersions                  []upgradeVersion      `json:"upgradeVersions"`
	UsedCapacity                     capacityLimit         `json:"usedCapacity"`
	UserName                         string                `json:"userName"`
	VscanFileOperationDefaultProfile string                `json:"vscanFileOperationDefaultProfile"`
	WormEnabled                      bool                  `json:"wormEnabled"`
	WritingSpeedState                string                `json:"writingSpeedState"`
}

type broadcastDomainInfo struct {
	BroadcastDomain string `json:"broadcastDomain"`
	IPSpace         string `json:"ipSpace"`
	Mtu             int    `json:"mtu"`
}

type capacityTierInfo struct {
	CapacityTierUsedSize capacityLimit `json:"capacityTierUsedSize"`
	S3BucketName         string        `json:"s3BucketName"`
	TierLevel            string        `json:"tierLeve"`
}

type node struct {
	CloudProviderID      string      `json:"cloudProviderId"`
	Health               bool        `json:"health"`
	InTakeover           bool        `json:"inTakeover"`
	Lifs                 []lif       `json:"lifs"`
	Name                 string      `json:"name"`
	PlatformLicense      interface{} `json:"platformLicense"`
	PlatformSerialNumber interface{} `json:"platformSerialNumber"`
	SerialNumber         string      `json:"serialNumber"`
	SystemID             string      `json:"systemId"`
}

type upgradeVersion struct {
	ImageVersion      string `json:"imageVersion"`
	LastModified      int    `json:"lastModified"`
	AutoUpdateAllowed bool   `json:"autoUpdateAllowed"`
}

type licenseType struct {
	CapacityLimit capacityLimit `json:"capacityLimit"`
	Name          string        `json:"name"`
}

type capacityLimit struct {
	Size float64 `json:"size"`
	Unit string  `json:"unit"`
}

type lif struct {
	DataProtocols []string `json:"dataProtocols"`
	IP            string   `json:"ip"`
	LifType       string   `json:"lifType"`
	Netmask       string   `json:"netmask"`
	NodeName      string   `json:"nodeName"`
	PrivateIP     bool     `json:"privateIp"`
}

type cvoStatus struct {
	ExtendedFailureReason interface{}   `json:"extendedFailureReason"`
	FailureCauses         failureCauses `json:"failureCauses"`
	Message               string        `json:"message"`
	Status                string        `json:"status"`
}

type failureCauses struct {
	InvalidCloudProviderCredentials bool `json:"invalidCloudProviderCredentials"`
	InvalidOntapCredentials         bool `json:"invalidOntapCredentials"`
	NoCloudProviderConnection       bool `json:"noCloudProviderConnection"`
}

// userTags the input for requesting a CVO
type userTags struct {
	TagKey   string `structs:"tagKey"`
	TagValue string `structs:"tagValue,omitempty"`
}

// modifyUserTagsRequest the input for requesting tags modificaiton
type modifyUserTagsRequest struct {
	Tags []userTags `structs:"tags"`
}

// setPasswordRequest the input for for setting password
type setPasswordRequest struct {
	Password string `structs:"password"`
}

// licenseAndInstanceTypeModificationRequest the input for license and instance type modification
type licenseAndInstanceTypeModificationRequest struct {
	InstanceType string `structs:"instanceType"`
	LicenseType  string `structs:"licenseType"`
}

// changeTierLevelRequest the input for tier level change
type changeTierLevelRequest struct {
	Level string `structs:"level"`
}

// upgradeOntapVersionRequest
type upgradeOntapVersionRequest struct {
	UpdateType      string `structs:"updateType"`
	UpdateParameter string `structs:"updateParameter"`
}

// set config flag
type setFlagRequest struct {
	Value     bool   `structs:"value"`
	ValueType string `structs:"valueType"`
}

// Check HTTP response code, return error if HTTP request is not successed.
func apiResponseChecker(statusCode int, response []byte, funcName string) error {

	if statusCode >= 300 || statusCode < 200 {
		log.Printf("%s request failed: %v", funcName, string(response))
		return fmt.Errorf("code: %d, message: %s", statusCode, string(response))
	}

	return nil

}

func (c *Client) checkTaskStatus(id string) (int, string, error) {

	log.Printf("checkTaskStatus: %s", id)

	baseURL := fmt.Sprintf("/occm/api/audit/activeTask/%s", id)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Printf("checkTaskStatus request failed: %v, %v", code, err)
				return 0, "", err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatus")
	if responseError != nil {
		return 0, "", responseError
	}

	var result cvoStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatus ", err)
		return 0, "", err
	}

	return result.Status, result.Error, nil
}

func (c *Client) waitOnCompletion(id string, actionName string, task string, retries int, waitInterval int) error {
	for {
		cvoStatus, failureErrorMessage, err := c.checkTaskStatus(id)
		if err != nil {
			return err
		}
		if cvoStatus == 1 {
			return nil
		} else if cvoStatus == -1 {
			return fmt.Errorf("Failed to %s %s, error: %s", task, actionName, failureErrorMessage)
		} else if cvoStatus == 0 {
			if retries == 0 {
				log.Print("Taking too long to ", task, actionName)
				return fmt.Errorf("Taking too long for %s to %s or not properly setup", actionName, task)
			}
			log.Printf("Sleep for %d seconds", waitInterval)
			time.Sleep(time.Duration(waitInterval) * time.Second)
			retries--
		}

	}
}

// get working environment information by working environment id
// response: publicId, name, isHA, cloudProvider, workingEnvironmentType
func (c *Client) getWorkingEnvironmentInfo(id string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/occm/api/working-environments/%s", id)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Printf("getWorkingEnvironmentInfo: ID %s request failed. Err: %v", id, err)
				return workingEnvironmentInfo{}, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentInfo")
	if responseError != nil {
		log.Printf("apiResponseChecker error %v", responseError)
		return workingEnvironmentInfo{}, responseError
	}

	var result workingEnvironmentInfo
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}

	return result, nil
}

func findWE(name string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].Name == name {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment %s in the list", name)
}

func findWEForID(id string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].PublicID == id {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment %s in the list", id)
}

func (c *Client) findWorkingEnvironmentByName(name string) (workingEnvironmentInfo, error) {
	// check working environment exists or not
	baseURL := fmt.Sprintf("/occm/api/working-environments/exists/%s", name)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("findWorkingEnvironmentByName request failed. (check exists) ", statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	// get working environment information
	baseURL = "/occm/api/working-environments"
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("findWorkingEnvironmentByName %s request failed (%d)", name, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError = apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentByName")
		return workingEnvironmentInfo{}, err
	}

	var workingEnvironment workingEnvironmentInfo
	workingEnvironment, err = findWE(name, workingEnvironments.VsaWorkingEnvironment)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.OnPremWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.AzureVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.GcpVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}

	log.Printf("Cannot find the working environment %s", name)

	return workingEnvironmentInfo{}, err
}

// get WE directly from REST API using a given ID
func (c *Client) findWorkingEnvironmentByID(id string) (workingEnvironmentInfo, error) {

	workingEnvInfo, err := c.getWorkingEnvironmentInfo(id)
	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_id %s", id)
	}
	workingEnvDetail, err := c.findWorkingEnvironmentByName(workingEnvInfo.Name)
	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_name %s", workingEnvInfo.Name)
	}
	return workingEnvDetail, nil
}

func (c *Client) getFSXWorkingEnvironmentInfo(tenantID string, id string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s", tenantID, id)
	hostType := "CloudManagerHost"
	var result workingEnvironmentInfo

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("getFSXWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		log.Printf("error: %#v", err)
		return workingEnvironmentInfo{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getFSXWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var system map[string]interface{}
	if err := json.Unmarshal(response, &system); err != nil {
		log.Print("Failed to unmarshall response from getFSXWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}
	result.Name = system["name"].(string)

	baseURL = fmt.Sprintf("/occm/api/fsx/working-environments/%s/svms", id)
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("getFSXWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}
	responseError = apiResponseChecker(statusCode, response, "getFSXWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}
	var info []map[string]interface{}
	if err := json.Unmarshal(response, &info); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}
	//assume there is only one svm in fsx
	result.SvmName = info[0]["name"].(string)

	return result, nil
}

func (c *Client) getAPIRoot(workingEnvironmentID string) (string, string, error) {

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", "", err
		}
		c.Token = accesTokenResult.Token
	}

	// fsx working environment starts with "fs-" prefix.
	if strings.HasPrefix(workingEnvironmentID, "fs-") {
		return "/occm/api/fsx", "", nil
	}
	workingEnvDetail, err := c.getWorkingEnvironmentInfo(workingEnvironmentID)
	if err != nil {
		log.Print("Cannot get working environment information.")
		return "", "", err
	}
	log.Printf("Working environment %v", workingEnvDetail)

	var baseURL string
	if workingEnvDetail.CloudProviderName != "Amazon" {
		if workingEnvDetail.IsHA {
			baseURL = fmt.Sprintf("/occm/api/%s/ha", strings.ToLower(workingEnvDetail.CloudProviderName))
		} else {
			baseURL = fmt.Sprintf("/occm/api/%s/vsa", strings.ToLower(workingEnvDetail.CloudProviderName))
		}
	} else {
		if workingEnvDetail.IsHA {
			baseURL = "/occm/api/aws/ha"
		} else {
			baseURL = "/occm/api/vsa"
		}
	}
	log.Printf("API root = %s", baseURL)
	return baseURL, workingEnvDetail.CloudProviderName, nil
}

func (c *Client) getAPIRootForWorkingEnvironment(isHA bool, workingEnvironmentID string) string {

	var baseURL string

	if workingEnvironmentID == "" {
		if isHA == true {
			baseURL = "/occm/api/gcp/ha/working-environments"
		} else {
			baseURL = "/occm/api/gcp/vsa/working-environments"
		}
	} else {
		if isHA == true {
			baseURL = fmt.Sprintf("/occm/api/gcp/ha/working-environments/%s", workingEnvironmentID)
		} else {
			baseURL = fmt.Sprintf("/occm/api/gcp/vsa/working-environments/%s", workingEnvironmentID)
		}
	}

	log.Printf("API root = %s", baseURL)
	return baseURL
}

// read working environemnt information and return the details
func (c *Client) getWorkingEnvironmentDetail(d *schema.ResourceData) (workingEnvironmentInfo, error) {
	var workingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("file_system_id"); ok {
		workingEnvDetail, err = c.getFSXWorkingEnvironmentInfo(d.Get("tenant_id").(string), a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_id %s", a.(string))
		}
		return workingEnvDetail, nil
	}

	if a, ok := d.GetOk("working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		workingEnvDetail, err = c.findWorkingEnvironmentByID(WorkingEnvironmentID)
		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_id %s", WorkingEnvironmentID)
		}
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", workingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}
	return workingEnvDetail, nil
}

func (c *Client) getFSXSVM(id string) (string, error) {

	log.Print("getFSXSVM")

	baseURL := fmt.Sprintf("/occm/api/fsx/working-environments/%s/svms", id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getFSXSVM request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getFSXSVM")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxSVMResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getFSXSVM ", err)
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no SVM found for %s", id)
	}

	return result[0].Name, nil
}

func (c *Client) getAWSFSXByName(name string, tenantID string) (string, error) {

	log.Print("getAWSFSXByName")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getAWSFSXByName request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getAWSFSXByName request failed ", statusCode, err)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getAWSFSXByName")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAWSFSXByName ", err)
		return "", err
	}

	for _, fsxID := range result {
		if fsxID.Name == name {
			return fsxID.ID, nil
		}
	}

	return "", nil
}

// read working environemnt information and return the details
func (c *Client) getWorkingEnvironmentDetailForSnapMirror(d *schema.ResourceData) (workingEnvironmentInfo, workingEnvironmentInfo, error) {
	var sourceWorkingEnvDetail workingEnvironmentInfo
	var destWorkingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("source_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		sourceWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID)
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_id %s", WorkingEnvironmentID)
		}
	} else if a, ok = d.GetOk("source_working_environment_name"); ok {
		sourceWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", sourceWorkingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_id or source_working_environment_name")
	}

	if a, ok := d.GetOk("destination_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		// fsx working environment starts with "fs-" prefix.
		if strings.HasPrefix(WorkingEnvironmentID, "fs-") {
			if b, ok := d.GetOk("tenant_id"); ok {
				tenantID := b.(string)
				_, err := c.getAWSFSX(WorkingEnvironmentID, tenantID)
				if err != nil {
					log.Print("Error getting AWS FSX")
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				destWorkingEnvDetail.PublicID = WorkingEnvironmentID
				svmName, err := c.getFSXSVM(WorkingEnvironmentID)
				if err != nil {
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				destWorkingEnvDetail.SvmName = svmName
			} else {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find FSX working environment by destination_working_environment_id %s, need tenant_id", WorkingEnvironmentID)
			}
		} else {
			destWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID)
			if err != nil {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_id %s", WorkingEnvironmentID)
			}
			log.Print("findWorkingEnvironmentForID", destWorkingEnvDetail)
		}
	} else if a, ok = d.GetOk("destination_working_environment_name"); ok {
		if b, ok := d.GetOk("tenant_id"); ok {
			workingEnvironmentName := a.(string)
			tenantID := b.(string)
			WorkingEnvironmentID, err := c.getAWSFSXByName(workingEnvironmentName, tenantID)
			if err != nil {
				log.Print("Error getting AWS FSX: ", err)
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
			}
			destWorkingEnvDetail.PublicID = WorkingEnvironmentID
			svmName, err := c.getFSXSVM(WorkingEnvironmentID)
			if err != nil {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
			}
			destWorkingEnvDetail.SvmName = svmName
		} else {
			destWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
			if err != nil {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_name %s", a.(string))
			}
			log.Printf("Get environment id %v by %v", destWorkingEnvDetail.PublicID, a.(string))
		}
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_id or destination_working_environment_name")
	}
	return sourceWorkingEnvDetail, destWorkingEnvDetail, nil
}

// get all WE from REST API and then using a given ID get the WE
func (c *Client) findWorkingEnvironmentForID(id string) (workingEnvironmentInfo, error) {
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	baseURL := "/occm/api/working-environments"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("findWorkingEnvironmentForId %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentForId")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentForId")
		return workingEnvironmentInfo{}, err
	}

	var workingEnvironment workingEnvironmentInfo
	workingEnvironment, err = findWEForID(id, workingEnvironments.VsaWorkingEnvironment)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.OnPremWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.AzureVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.GcpVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}

	log.Printf("Cannot find the working environment %s", id)

	return workingEnvironmentInfo{}, err
}

// get working environment properties
func (c *Client) getWorkingEnvironmentProperties(apiRoot string, id string, field string) (workingEnvironmentOntapClusterPropertiesResponse, error) {
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("%s/working-environments/%s?fields=%s", apiRoot, id, field)
	log.Printf("Call %s", baseURL)

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Printf("getWorkingEnvironmentProperties %s request failed (%d) %s", baseURL, statusCode, err)
				return workingEnvironmentOntapClusterPropertiesResponse{}, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentProperties")
	if responseError != nil {
		return workingEnvironmentOntapClusterPropertiesResponse{}, responseError
	}

	var result workingEnvironmentOntapClusterPropertiesResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentProperties ", err)
		return workingEnvironmentOntapClusterPropertiesResponse{}, err
	}
	log.Printf("Get cvo properities result %+v", result)
	return result, nil
}

// customized check diff user-tags (aws_tag, azure_tag, gcp_label)
func checkUserTagDiff(diff *schema.ResourceDiff, tagName string, keyName string) error {
	if diff.HasChange(tagName) {
		_, expectTag := diff.GetChange(tagName)
		etags := expectTag.(*schema.Set)
		if etags.Len() > 0 {
			log.Println("etags len: ", etags.Len())
			// check each of the tag_key in the list is unique
			respErr := checkUserTagKeyUnique(etags, keyName)
			if respErr != nil {
				return respErr
			}
		}
	}
	return nil
}

// check each of the tag_key or label_key in the list is unique
func checkUserTagKeyUnique(etags *schema.Set, keyName string) error {
	m := make(map[string]bool)
	for _, v := range etags.List() {
		tag := v.(map[string]interface{})
		tkey := tag[keyName].(string)
		if _, ok := m[tkey]; !ok {
			m[tkey] = true
		} else {
			return fmt.Errorf("%s %s is not unique", keyName, tkey)
		}
	}
	return nil
}

// expandUserTags converts set to userTags struct
func expandUserTags(set *schema.Set) []userTags {
	tags := []userTags{}

	for _, v := range set.List() {
		tag := v.(map[string]interface{})
		userTag := userTags{}
		userTag.TagKey = tag["tag_key"].(string)
		userTag.TagValue = tag["tag_value"].(string)
		tags = append(tags, userTag)
	}
	return tags
}

func (c *Client) callCMUpdateAPI(method string, request interface{}, baseURL string, id string, functionName string) error {
	apiRoot, _, err := c.getAPIRoot(id)
	baseURL = apiRoot + baseURL

	hostType := "CloudManagerHost"
	params := structs.Map(request)

	if c.Token == "" {
		accessTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Printf("in %s request, failed to get AccessToken", functionName)
			return err
		}
		c.Token = accessTokenResult.Token
	}

	statusCode, response, _, err := c.CallAPIMethod(method, baseURL, params, c.Token, hostType)
	if err != nil {
		log.Printf("%s request failed: %d", functionName, statusCode)
		log.Print("call api response: ", response)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, functionName)
	if responseError != nil {
		return responseError
	}
	return nil
}

// update CVO user-tags
func updateCVOUserTags(d *schema.ResourceData, meta interface{}, tagName string) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request modifyUserTagsRequest
	if c, ok := d.GetOk(tagName); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			if tagName == "gcp_label" {
				request.Tags = expandGCPLabelsToUserTags(tags)
			} else {
				request.Tags = expandUserTags(tags)
			}
			log.Print("Update user-tags: ", request.Tags)
		}
	}
	// Update tags
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/user-tags", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOUserTags")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s %s: %v", id, tagName, request.Tags)
	return nil
}

// set the cluster password of a specific cloud volumes ONTAP
func updateCVOSVMPassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request setPasswordRequest
	request.Password = d.Get("svm_password").(string)
	// Update password
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/set-password", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOSVMPassword")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s svm_password", id)
	return nil
}

func (c *Client) waitOnCompletionCVOUpdate(apiRoot string, id string, retryCount int, waitInterval int) error {
	// check upgrade status
	log.Print("Check CVO update status")

	for {
		cvoResp, err := c.getWorkingEnvironmentProperties(apiRoot, id, "status,ontapClusterProperties")
		if err != nil {
			return err
		}
		if cvoResp.Status.Status != "UPDATING" {
			log.Print("CVO update is done")
			return nil
		}
		if retryCount <= 0 {
			log.Print("Taking too long for status to be active")
			return fmt.Errorf("Taking too long for CVO to be active or not properly setup")
		}
		log.Printf("Update status %s...(%d)", cvoResp.Status.Status, retryCount)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retryCount--
	}
}

// set the license_type and instance type of a specific cloud volumes ONTAP
func updateCVOLicenseInstanceType(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request licenseAndInstanceTypeModificationRequest
	if c, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = c.(string)
	}
	if c, ok := d.GetOk("license_type"); ok {
		request.LicenseType = c.(string)
	}

	// Update license type and instance type
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/license-instance-type", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOLicenseInstanceType")
	if updateErr != nil {
		return updateErr
	}
	// check upgrade status
	apiRoot, _, err := client.getAPIRoot(id)
	if err != nil {
		return fmt.Errorf("Cannot get root API")
	}

	retryCount := 65
	if d.Get("is_ha").(bool) {
		retryCount = retryCount * 2
	}
	err = client.waitOnCompletionCVOUpdate(apiRoot, id, retryCount, 60)
	if err != nil {
		return fmt.Errorf("Update CVO failed %v", err)
	}
	log.Printf("Updated %s license and instance type: %v", id, request)
	return nil
}

// update tier_level of a specific cloud volumes ONTAP
func updateCVOTierLevel(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request changeTierLevelRequest
	if c, ok := d.GetOk("tier_level"); ok {
		request.Level = c.(string)
	}
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/change-tier-level", id)
	updateErr := client.callCMUpdateAPI("POST", request, baseURL, id, "updateCVOTierLevel")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s tier_level: %v", id, request)
	return nil
}

func (c *Client) waitOnCompletionOntapImageUpgrade(apiRoot string, id string, targetVersion string, retryCount int, waitInterval int) error {
	// check upgrade status
	log.Print("Check CVO ontap image upgrade status")

	for {
		cvoResp, err := c.getWorkingEnvironmentProperties(apiRoot, id, "status,ontapClusterProperties")
		if err != nil {
			return err
		}
		if cvoResp.Status.Status != "UPDATING" && cvoResp.OntapClusterProperties.OntapVersion != "" {
			if strings.Contains(targetVersion, cvoResp.OntapClusterProperties.OntapVersion) {
				log.Print("ONTAP image upgrade is done")
				return nil
			}
			log.Printf("Update ontap image failed on checking version (%s, %s)", cvoResp.OntapClusterProperties.OntapVersion, targetVersion)
			return fmt.Errorf("Update ontap version failed. Current version %s", cvoResp.OntapClusterProperties.OntapVersion)
		}
		if retryCount <= 0 {
			log.Print("Taking too long for status to be active")
			return fmt.Errorf("Taking too long for CVO to be active or not properly setup")
		}
		log.Printf("Update %s status %s...(%d)", targetVersion, cvoResp.Status.Status, retryCount)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retryCount--
	}
}

// check if ontap_version is the list of upgrade available versions
func (c *Client) upgradeOntapVersionAvailable(apiRoot string, id string, ontapVersion string) (string, error) {
	log.Print("upgradeOntapVersionAvailable: Check if target version is in the upgrade version list")

	var upgradeOntapVersions []upgradeVersion

	WEProperties, err := c.getWorkingEnvironmentProperties(apiRoot, id, "ontapClusterProperties.fields(upgradeVersions)")
	if err != nil {
		return "", fmt.Errorf("upgradeOntapVersionAvailable %s not able to get the properties %v", id, err)
	}
	log.Printf("Get current ontap version: %s", WEProperties.OntapClusterProperties.OntapVersion)

	upgradeOntapVersions = WEProperties.OntapClusterProperties.UpgradeVersions

	if upgradeOntapVersions != nil {
		for _, ugVersion := range upgradeOntapVersions {
			version := ugVersion.ImageVersion
			if strings.Contains(ontapVersion, version) {
				return version, nil
			}
		}
		return "", fmt.Errorf("Working environment %s: ontap version %s is not in the upgrade versions list (%+v)", id, ontapVersion, upgradeOntapVersions)
	}
	return "", fmt.Errorf("Working environment %s: no upgrade version availble", id)
}

func (c *Client) setConfigFlag(request setFlagRequest, keyPath string) error {
	log.Print("setConfigFlag: set flag to allow ONTAP image upgrade")

	hostType := "CloudManagerHost"

	baseURL := fmt.Sprintf("/occm/api/occm/config/%s", keyPath)
	params := structs.Map(request)
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, params, c.Token, hostType)

	responseError := apiResponseChecker(statusCode, response, "setUpgradeCheckingBypass")
	if responseError != nil {
		return responseError
	}

	if err != nil {
		log.Print("setUpgradeCheckingBypass request failed ", statusCode)
		return err
	}
	return nil

}

// upgrade CVO ontap version
func (c *Client) upgradeCVOOntapImage(apiRoot string, id string, ontapVersion string, isHa bool) error {
	// set config flag to skip the upgrade check
	var setFlag setFlagRequest
	setFlag.Value = true
	setFlag.ValueType = "BOOLEAN"

	log.Printf("Set config flag")
	setFlagErr := c.setConfigFlag(setFlag, "skip-eligibility-paygo-upgrade")
	if setFlagErr != nil {
		log.Printf("upgradeCVOOntapVersion failed on setConfigFlag call %v", setFlagErr)
		return setFlagErr
	}

	// upgrade image
	var request upgradeOntapVersionRequest
	request.UpdateType = "OCCM_PROVIDED"
	request.UpdateParameter = ontapVersion

	baseURL := fmt.Sprintf("/working-environments/%s/update-image", id)
	log.Printf("upgradeCVOOntapVersion - %s %v", baseURL, request)
	updateErr := c.callCMUpdateAPI("POST", request, baseURL, id, "upgradeCVOOntapVersion")
	if updateErr != nil {
		log.Printf("upgradeCVOOntapVersion failed on API call %v", updateErr)
		return updateErr
	}

	// check upgrade status
	retryCount := 65
	if isHa {
		retryCount = retryCount * 2
	}
	err := c.waitOnCompletionOntapImageUpgrade(apiRoot, id, ontapVersion, retryCount, 60)
	if err != nil {
		return fmt.Errorf("Upgrade ontap image %s failed %v", ontapVersion, err)
	}
	log.Printf("Updated %s ontap_version: %v", id, request)
	return nil
}

func (c *Client) doUpgradeCVOOntapVersion(id string, isHA bool, ontapVersion string) error {
	// only when the upgrade_ontap_version is true, use_latest_version is false and the ontap_version is not "latest"
	log.Print("Check CVO ontap image upgrade status ... ")
	apiRoot, _, err := c.getAPIRoot(id)
	if err != nil {
		return fmt.Errorf("Cannot get root API")
	}

	upgradeVersion, err := c.upgradeOntapVersionAvailable(apiRoot, id, ontapVersion)
	if err != nil {
		return err
	}

	return c.upgradeCVOOntapImage(apiRoot, id, upgradeVersion, isHA)
}

func checkOntapVersionChangeWithoutUpgrade(d *schema.ResourceData) error {
	var wrongChange = false
	if d.HasChange("ontap_version") {
		currentVersion, _ := d.GetChange("ontap_version")
		d.Set("ontap_version", currentVersion)
		wrongChange = true
	}
	if d.HasChange("use_latest_version") {
		current, _ := d.GetChange("use_latest_version")
		d.Set("use_latest_version", current)
		wrongChange = true
	}
	if wrongChange {
		return fmt.Errorf("upgrade_ontap_version is not turned on. The change will not be done")
	}
	log.Printf("No ontap version upgrade")
	return nil
}

func (c *Client) checkAndDoUpgradeOntapVersion(d *schema.ResourceData) error {
	upgradeOntapVersion := d.Get("upgrade_ontap_version").(bool)
	if upgradeOntapVersion {
		ontapVersion := d.Get("ontap_version").(string)
		log.Printf("Check if need upgrade - ontapVersion %s", ontapVersion)

		if ontapVersion == "latest" {
			return fmt.Errorf("ontap_version only can be upgraded with the specific ontap_version not \"latest\"")
		}
		if d.Get("use_latest_version").(bool) {
			return fmt.Errorf("ontap_version cannot be upgraded with \"use_latest_version\" true")
		}
		id := d.Id()
		respErr := c.doUpgradeCVOOntapVersion(id, d.Get("is_ha").(bool), ontapVersion)
		if respErr != nil {
			currentVersion, _ := d.GetChange("ontap_version")
			d.Set("ontap_version", currentVersion)
			return respErr
		}
	} else {
		respErr := checkOntapVersionChangeWithoutUpgrade(d)
		if respErr != nil {
			return respErr
		}
	}
	return nil
}
