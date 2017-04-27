package channels

import (
	"accountcreationautomation/constants"
	"accountcreationautomation/logger"
	"accountcreationautomation/model"
	"accountcreationautomation/post"
	"accountcreationautomation/stack"
	"encoding/csv"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var mu sync.RWMutex
var wg sync.WaitGroup

func New(records [][]string, ApiAccessKey, ApiSecretKey, DomainUrl string, writer *csv.Writer) {
	lenArray := len(records)
	//Creating a Buffered Channel
	inputChan := make(chan model.CsvWriter, lenArray)

	for i, record := range records {

		accesskey := record[0]
		secretkey := record[1]
		region := record[2]
		AccountName := record[3]

		// skip header row
		if i == 0 {
			writer.Write([]string{accesskey, secretkey, region, AccountName, constants.DOMAIN, constants.CSV_HEADER_ROLEARN, constants.CSV_HEADER_STATUS, constants.CSV_HEADER_CMPUTE_ACCOUNT})
			continue
		}
		wg.Add(1)

		token := ""

		defer wg.Wait()

		go func() {
			StartAccountCreation(accesskey, secretkey, region, AccountName, DomainUrl, token, ApiAccessKey, ApiSecretKey, inputChan, writer)
			FileWriterFunc(inputChan, writer)
		}()

	}

}

func StartAccountCreation(accesskey, secretkey, region, AccountName, DomainUrl, token, ApiAccessKey, ApiSecretKey string, inputChan chan model.CsvWriter, writer *csv.Writer) {
	time.Sleep(1 * time.Second)
	defer wg.Done()
	creds := credentials.NewStaticCredentials(accesskey, secretkey, token)
	//Creating CloudFormationStack in us-east-1
	svc := cloudformation.New(session.New(&aws.Config{Region: aws.String(constants.REGION), Credentials: creds}))
	var stackid string

	//Creating CloudFormationStack using Batchly's CloudFormationTemplate
	createstackResp, err := stack.CreateStack(svc)
	if err != nil {
		mu.Lock()
		defer mu.Unlock()
		logger.Get().Error(err)
		time.Sleep(1 * time.Second)
		CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: constants.ROLEARN_ERROR, RoleArnStatus: err.Error(), CmputeAccountStatus: constants.CMPUTE_ACCOUNT_ERROR}
		val := CsvWriterObj
		writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
		inputChan <- CsvWriterObj

	} else {
		stackid = *createstackResp.StackId
		//Checking the status of the stack,loop exits when the status is either Success or Failure
	describestatus:
		describestackResp, err := stack.DescribeStack(svc, stackid)
		if err != nil {
			fmt.Println(err)
			logger.Get().Error("Describe Stack Error ", err)
		}

		for _, stack := range describestackResp.Stacks {
			if *stack.StackStatus == "CREATE_FAILED" {
				mu.Lock()
				defer mu.Unlock()
				CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: constants.ROLEARN_ERROR, RoleArnStatus: constants.ROLEARN_FAILURE_STATUS, CmputeAccountStatus: constants.CMPUTE_ACCOUNT_ERROR}
				val := CsvWriterObj
				writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
				inputChan <- CsvWriterObj
			} else if *stack.StackStatus == "CREATE_COMPLETE" {
				for _, output := range stack.Outputs {

					Result := post.PostCall(AccountName, region, *output.OutputValue, "ApiKey "+ApiAccessKey+":"+ApiSecretKey, DomainUrl, 0)
					if Result == 1 {
						mu.Lock()
						defer mu.Unlock()
						CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: *output.OutputValue, RoleArnStatus: constants.ROLEARN_SUCCESS_STATUS, CmputeAccountStatus: constants.CMPUTE_ACCOUNT}
						val := CsvWriterObj
						writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
						inputChan <- CsvWriterObj

					} else {
						mu.Lock()
						defer mu.Unlock()
						CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: *output.OutputValue, RoleArnStatus: constants.ROLEARN_SUCCESS_STATUS, CmputeAccountStatus: constants.CMPUTE_ACCOUNT_ERROR}
						val := CsvWriterObj
						writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
						inputChan <- CsvWriterObj

					}
				}
			} else if *stack.StackStatus == "CREATE_IN_PROGRESS" {
				time.Sleep(3 * time.Second)
				goto describestatus

			} else if *stack.StackStatus == "ROLLBACK_IN_PROGRESS" {
				time.Sleep(3 * time.Second)
				goto describestatus

			} else if *stack.StackStatus == "ROLLBACK_FAILED" {
				mu.Lock()
				defer mu.Unlock()
				CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: constants.ROLEARN_ERROR, RoleArnStatus: constants.ROLEARN_FAILURE_STATUS, CmputeAccountStatus: constants.CMPUTE_ACCOUNT_ERROR}
				val := CsvWriterObj
				writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
				inputChan <- CsvWriterObj

			} else if *stack.StackStatus == "ROLLBACK_COMPLETE" {
				mu.Lock()
				defer mu.Unlock()
				CsvWriterObj := model.CsvWriter{Accesskey: accesskey, Secretkey: secretkey, Region: region, Account: AccountName, Domain: DomainUrl, RoleArn: constants.ROLEARN_ERROR, RoleArnStatus: constants.ROLEARN_FAILURE_STATUS, CmputeAccountStatus: constants.CMPUTE_ACCOUNT_ERROR}
				val := CsvWriterObj
				writer.Write([]string{val.Accesskey, val.Secretkey, val.Region, val.Account, val.Domain, val.RoleArn, val.RoleArnStatus, val.CmputeAccountStatus})
				inputChan <- CsvWriterObj

			}
		}

	}
}

//Getting the data from the Channel
func FileWriterFunc(inputChan chan model.CsvWriter, writer *csv.Writer) {
	//	defer close(inputChan)

	for val := range inputChan {
		logger.Get().Infoln("Getting the Value out of channel", val)
	}

}
