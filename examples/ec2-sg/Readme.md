## Example
This example creates a security group stack, an ec2 instance stack. It contains compose configuration file, security group template and ec2 instance cfn template.

NOTE: Please update the `VpcId` and `SubnetId`

## Configuration File
Compose File
```yaml
Description: Sample CloudFormation Compose file
Vars:
  ENV_NAME: cfnc
  ENV_TYPE: nonproduction
  SUBNET_ID: "subnet-033274e18559e7bde"
  VPC_ID: "vpc-0001e3b703212c9cb"
Flows:
  SecurityGroup:
    Order: 0
    Description: Creates Security Group
    Stacks:
    - template_file: sg.yml
      stack_name: sample-{{ .ENV_NAME }}-security-group
      parameters:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
        VpcId: '{{ .VPC_ID }}'
      tags:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'

  EC2Instance:
    Order: 1
    Description: Creates EC2 Instance
    Stacks:
       - template_file: ec2.yml
      stack_name: sample-{{ .ENV_NAME }}-ec2-instance
      parameters:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
        SubnetId: '{{ .SUBNET_ID }}'
      tags:
        EnvironmentName: '{{ .ENV_NAME }}'
        EnvironmentType: '{{ .ENV_TYPE }}'
```
Security Group CFN template file: sg.yml
```yaml
Parameters:
  EnvironmentType:
    Type: String
  EnvironmentName:
    Type: String
  VpcId:
    Type: String

Resources:
  CfnComposeSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: !Sub "CfnCompose test security group for ${EnvironmentType}"
      VpcId: !Ref VpcId
      Tags:
        - Key: EnvironmentName
          Value: !Ref EnvironmentName
        - Key: EnvironmentType
          Value: !Ref EnvironmentType

Outputs:
  SecurityGroupId:
    Value: !Ref CfnComposeSecurityGroup
    Export:
      Name: !Sub ${EnvironmentName}:CfnComposeSecurityGroupId
```

EC2 Instance CFN template file: ec2.yaml
```yaml
Parameters:
  EnvironmentType:
    Type: String
  EnvironmentName:
    Type: String
  SubnetId:
    Type: String
  ImageId:
    Type: AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>
    Default: '/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2'
  InstanceType:
    Type: String
    Default: t2.nano
Resources:
  ExampleEC2Instance:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: !Ref ImageId
      InstanceType: !Ref InstanceType
      SubnetId: !Ref SubnetId
      SecurityGroupIds:
        - Fn::ImportValue: !Sub ${EnvironmentName}:CfnComposeSecurityGroupId
      Tags:
        - Key: EnvironmentName
          Value: !Ref EnvironmentName
        - Key: EnvironmentType
          Value: !Ref EnvironmentType
```
## Visualize
`cfnc config visualize`

## Deploy
`cfnc deploy`

## Destroy
`cfnc destroy`
