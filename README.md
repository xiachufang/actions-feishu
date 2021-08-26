# actions-feishu

[![Lint](https://github.com/xiachufang/actions-feishu/actions/workflows/lint.yml/badge.svg)](https://github.com/xiachufang/actions-feishu/actions/workflows/lint.yml)
[![build-test](https://github.com/xiachufang/actions-feishu/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/xiachufang/actions-feishu/actions/workflows/test.yml)

通过 GitHub Actions 来发送消息到飞书

# Quick start

Actions 配置样例：

```yaml
    - name: Send feishu message
      env:
        ACTIONS_FEISHU_TAG: 'v1.3.1' # 替换此变量, 最新见 https://github.com/xiachufang/actions-feishu/releases
        INPUT_WEBHOOK: "${{ secrets.FEISHU_ROBOT_WEBHOOK_URL }}"
        INPUT_TITLE: "I'm title"
        INPUT_CONTENT: "I'm message body\nfrom: ${{ github.repository }}"
      run: |
        wget -q https://github.com/xiachufang/actions-feishu/releases/download/${{ env.ACTIONS_FEISHU_TAG }}/linux-amd64-actions-feishu.tar.gz
        tar zxf linux-amd64-actions-feishu.tar.gz feishu
        ./feishu
```

更多示范例子见: [test.yml](./.github/workflows/test.yml)

# Configuration

## Inputs

| Variable | Required | Description |
| :---: | :---: | :----: |
| `webhook`| **true** | webhook address |
| `title` | **false** | title of message|
| `content` | **true** | content of message|
| `message_type`| **false**| message type, optional: `post`, `text`, `template`, default: `post`|
| `msg_template_path`| **false**| message template path, only work when `message_type` is `template`|
| `msg_template_values`| **false**| message template fillings, , only work when `message_type` is `template`|

详细定义见: [action.yml](./action.yml)

## Outputs

| Variable  | Description |
| :---:  | :----: |
| `response` | API response from feishu |
