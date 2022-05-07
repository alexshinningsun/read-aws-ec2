package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type cred struct {
	user string //AWS_ACCESS_KEY_ID
	pass string //AWS_SECRET_ACCESS_KEY
}

func main() {
	cre := &cred{
		user: "AWS_ACCESS_KEY_ID",
		pass: "AWS_SECRET_ACCESS_KEY",
	}
	var ec2Svc *ec2.EC2
	var r0 *ec2.DescribeRegionsOutput
	var r1 *ec2.DescribeInstancesOutput
	var r2 *ec2.DescribeSecurityGroupsOutput
	var r3 *ec2.DescribeKeyPairsOutput

	sess, err := newSess(cre, endpoints.ApEast1RegionID) // hard coded to Ap-East-1 Region, more detail plz check https://docs.aws.amazon.com/sdk-for-go/api/aws/endpoints/
	if err != nil {
		goto responseError
	}
	ec2Svc = ec2.New(sess)
	r0, err = ec2Svc.DescribeRegions(nil)
	if err != nil {
		goto responseError
	}
	fmt.Println(r0)

	r1, err = ec2Svc.DescribeInstances(nil)
	if err != nil {
		goto responseError
	}
	fmt.Println(r1)
	r2, err = ec2Svc.DescribeSecurityGroups(nil)
	if err != nil {
		goto responseError
	}
	fmt.Println(r2)
	r3, err = ec2Svc.DescribeKeyPairs(nil)
	if err != nil {
		goto responseError
	}
	fmt.Println(r3)
	return
responseError:
	fmt.Println(err)
	return
}

func newSess(c *cred, reg string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.user, c.pass, ""),
		Region:      aws.String(reg),
	})
	return sess, err
}
