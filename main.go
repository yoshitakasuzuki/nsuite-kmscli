package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/doublejumptokyo/nsuite-kmscli/awseoa"
	"github.com/doublejumptokyo/nsuite-kmscli/kmsutil"
)

var (
	flagTags = true
	prefix   = ""
)

func List(svc *kms.Client) (err error) {

	in := &kms.ListAliasesInput{}
	out, err := svc.ListAliases(context.TODO(), in)
	if err != nil {
		return
	}

	for _, a := range out.Aliases {
		alias := "None"
		if a.AliasName != nil {
			alias = *a.AliasName
		}
		alias = strings.TrimPrefix(alias, "alias/")
		if strings.HasPrefix(alias, "aws/") {
			continue
		}
		keyID := "None"
		if a.TargetKeyId != nil {
			keyID = *a.TargetKeyId
		}

		tags := ""
		if flagTags {
			in := &kms.ListResourceTagsInput{KeyId: a.TargetKeyId}
			out, err := svc.ListResourceTags(context.TODO(), in)
			if err != nil {
				return err
			}

			for _, t := range out.Tags {
				tags += *t.TagKey + ":" + *t.TagValue + "\t"
			}
		}

		fmt.Println(alias, keyID, tags)
	}
	return
}

func AddTag(svc *kms.Client, keyID, tagKey, tagValue string) (err error) {
	in := &kms.TagResourceInput{
		KeyId: aws.String(keyID),
		Tags: []kmstypes.Tag{
			{
				TagKey:   aws.String(tagKey),
				TagValue: aws.String(tagValue),
			},
		},
	}
	_, err = svc.TagResource(context.TODO(), in)
	return
}

func New(svc *kms.Client) (err error) {
	signer, err := awseoa.CreateSigner(svc, big.NewInt(4), prefix)
	if err != nil {
		return
	}
	fmt.Println(signer.Address().String(), signer.ID)
	return
}

func usage() {
	fmt.Println("Usage of awseoa:")
	fmt.Println("")
	fmt.Println("   list     Show list of keys")
	fmt.Println("            --tags: with tags")
	fmt.Println("   new      Create key")
	fmt.Println("            --prefix: alias prefix")
	fmt.Println("   add-tags [keyID] [name:value] [name:value]...")
	fmt.Println("            add tag to exist key")
}

func main() {
	var err error
	listFlag := flag.NewFlagSet("list", flag.ExitOnError)
	newFlag := flag.NewFlagSet("new", flag.ExitOnError)
	_ = flag.NewFlagSet("add-tags", flag.ExitOnError)

	listFlag.BoolVar(&flagTags, "tags", flagTags, "Show tags")
	newFlag.StringVar(&prefix, "prefix", "", "Prefix of alias")

	if len(os.Args) == 1 {
		usage()
		return
	}

	svc, err := kmsutil.NewKMSClient()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "list":
		if err := listFlag.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		err = List(svc)
	case "new":
		if err := newFlag.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		err = New(svc)
	case "add-tags":
		keyID := os.Args[2]

		for i := 3; i < len(os.Args); i++ {
			parts := strings.Split(os.Args[i], ":")
			err = AddTag(svc, keyID, parts[0], parts[1])
		}
	default:
		usage()
	}

	if err != nil {
		panic(err)
	}
}
