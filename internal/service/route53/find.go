package route53

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindHealthCheckByID(conn *route53.Route53, id string) (*route53.HealthCheck, error) {
	input := &route53.GetHealthCheckInput{
		HealthCheckId: aws.String(id),
	}

	output, err := conn.GetHealthCheck(input)

	if tfawserr.ErrCodeEquals(err, route53.ErrCodeNoSuchHealthCheck) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if output == nil || output.HealthCheck == nil || output.HealthCheck.HealthCheckConfig == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.HealthCheck, nil
}

func FindHostedZoneDNSSEC(conn *route53.Route53, hostedZoneID string) (*route53.GetDNSSECOutput, error) {
	input := &route53.GetDNSSECInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	output, err := conn.GetDNSSEC(input)

	if err != nil {
		return nil, err
	}

	return output, nil
}

func FindKeySigningKey(conn *route53.Route53, hostedZoneID string, name string) (*route53.KeySigningKey, error) {
	input := &route53.GetDNSSECInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	var result *route53.KeySigningKey

	output, err := conn.GetDNSSEC(input)

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, nil
	}

	for _, keySigningKey := range output.KeySigningKeys {
		if keySigningKey == nil {
			continue
		}

		if aws.StringValue(keySigningKey.Name) == name {
			result = keySigningKey
			break
		}
	}

	return result, err
}

func FindKeySigningKeyByResourceID(conn *route53.Route53, resourceID string) (*route53.KeySigningKey, error) {
	hostedZoneID, name, err := KeySigningKeyParseResourceID(resourceID)

	if err != nil {
		return nil, fmt.Errorf("error parsing Route 53 Key Signing Key (%s) identifier: %w", resourceID, err)
	}

	return FindKeySigningKey(conn, hostedZoneID, name)
}

func FindTrafficPolicyById(ctx context.Context, conn *route53.Route53, trafficPolicyId string) (*route53.TrafficPolicySummary, error) {
	var idMarker *string

	for allPoliciesListed := false; !allPoliciesListed; {
		input := &route53.ListTrafficPoliciesInput{}

		if idMarker != nil {
			input.TrafficPolicyIdMarker = idMarker
		}

		listResponse, err := conn.ListTrafficPoliciesWithContext(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, summary := range listResponse.TrafficPolicySummaries {
			if aws.StringValue(summary.Id) == trafficPolicyId {
				return summary, nil
			}
		}

		if aws.BoolValue(listResponse.IsTruncated) {
			idMarker = listResponse.TrafficPolicyIdMarker
		} else {
			allPoliciesListed = true
		}
	}

	return nil, nil
}

func FindTrafficPolicyInstanceId(ctx context.Context, conn *route53.Route53, id string) (*route53.GetTrafficPolicyInstanceOutput, error) {
	input := &route53.GetTrafficPolicyInstanceInput{
		Id: aws.String(id),
	}

	resp, err := conn.GetTrafficPolicyInstanceWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, route53.ErrCodeNoSuchTrafficPolicyInstance) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, &resource.NotFoundError{
			Message:     "Empty result",
			LastRequest: input,
		}
	}

	return resp, nil
}
