package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
    "github.com/aws/aws-sdk-go/aws/session"
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

func scan (ip string) {
	fmt.Printf("Scanning for %s\n", ip)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	svc := route53.New(sess)

	l, err := svc.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		exitErrorf("Unable to get hosted zones list, %v.", err)
	}
	for _, v := range l.HostedZones {
		// fmt.Printf("Scanning %s...\n", *v.Name)
		r, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{ HostedZoneId: v.Id })
		if err != nil {
			exitErrorf("Unable to get hosted zones list, %v.", err)
		}
		for _, resource := range r.ResourceRecordSets {
			if *resource.Type == "A" {
				for _, record := range resource.ResourceRecords {
					ipAddress := *record.Value
					if ipAddress == ip {
						fmt.Printf("%s %v\n", *resource.Name, *record.Value)
					}
				}
			}
		}
	}
}

func main() {
	scan(os.Args[1])
}