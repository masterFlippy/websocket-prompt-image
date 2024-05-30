import React from "react";
import { Button, Flex, Layout, Spin, Image, Switch } from "antd";
import TextArea from "antd/es/input/TextArea";
import { Typography } from "antd";
import { useWebSocket } from "../../contexts/websocket";

const { Paragraph, Text } = Typography;

const { Header, Content } = Layout;

const headerStyle: React.CSSProperties = {
  textAlign: "center",
  color: "#fff",
  height: 64,
  paddingInline: 48,
  lineHeight: "64px",
  backgroundColor: "#00415a",
};

const contentStyle: React.CSSProperties = {
  textAlign: "center",
  minHeight: 120,
  lineHeight: "120px",
  color: "#fff",
  backgroundColor: "#00719c",
};

const layoutStyle = {
  borderRadius: 8,
  overflow: "hidden",
  width: "calc(50% - 8px)",
  maxWidth: "calc(50% - 8px)",
  marginTop: "10vh",
  height: "100%",
};

const textAreaStyle = {
  width: "80%",
  height: "80%",
};

const buttonStyle = {
  width: "80%",
  height: "80%",
  marginBottom: 16,
  backgroundColor: "#009bd6",
};
const paragraphStyle = {
  margin: 10,
};

const textStyle = {
  color: "white",
};

export default function MainContent() {
  const id = "testing12asswwxc222";
  const { sendMessage, isLoading, url } = useWebSocket();
  const [bedrock, setBedrock] = React.useState(false);
  const [text, setText] = React.useState("");

  const handleBedrockChange = (checked: boolean) => {
    setBedrock(checked);
  };

  const handleSendMessage = () => {
    const message = {
      action: "generate",
      text: text,
      s3Key: id,
      bedrock,
    };
    sendMessage(message);
    setText("");
  };

  return (
    <Spin spinning={isLoading}>
      <Flex justify="center">
        <Layout style={layoutStyle}>
          <Header style={headerStyle}>Epic Post To Image POC</Header>
          <Content style={contentStyle}>
            <Paragraph style={paragraphStyle}>
              <Text strong style={textStyle}>
                Write a post as you would for a social media platform. Click
                generate to create an image for your post.
              </Text>
              .
            </Paragraph>

            <TextArea
              style={textAreaStyle}
              rows={4}
              placeholder="Write here"
              onChange={(e) => setText(e.target.value)}
            />
            <Paragraph style={paragraphStyle}>
              <Text strong style={textStyle}>
                Use Bedrock? (Toggle to use bedrock)
              </Text>

              <Switch onChange={handleBedrockChange} />
            </Paragraph>
            {url && (
              <Image width={"80%"} style={{ marginTop: "10px" }} src={url} />
            )}
            <Button
              type="primary"
              style={buttonStyle}
              onClick={handleSendMessage}
            >
              Generate
            </Button>
          </Content>
        </Layout>
      </Flex>
    </Spin>
  );
}
