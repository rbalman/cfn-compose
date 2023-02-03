### Demo Commands:
In this demo we will create **EC2Instance** and its required dependent stacks **SQS** and **RDS**. `SQS` and `RDS` wil be created first and then `EC2Instance` stack will be created.

https://user-images.githubusercontent.com/8892649/216241617-9aef4f3c-2981-4b36-a4da-41eec0e6b1e7.mp4/

- Deploy Stack in Dry Run Mode
  ```shell
    cfnc deploy -c cfnc.yml --dry-run
  ```

- Apply the real change and verify
    ```shell
      cfnc deploy -c cfnc.yml
    ```

- Destroy change in dry-run mode
  ```shell
    cfnc destroy -c cfnc.yml --dry-run
  ```

- Destroy
    ```shell
      cfnc destroy -c cfnc.yml
    ```