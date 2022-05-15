### find AZs of a region

```
aws --profile=default ec2 describe-availability-zones --region cn-northwest-1 --query 'AvailabilityZones[].ZoneName'
```
