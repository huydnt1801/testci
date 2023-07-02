import { createRef, useState, useEffect } from "react";
import {
    Modal,
    TouchableOpacity,
    View,
    Text,
    TouchableWithoutFeedback
} from "react-native";
import { useTranslation } from "react-i18next";

import className from "./className";

const messageDialogRef = createRef();

const MessageDialog = () => {

    const { t } = useTranslation();

    const [show, setShow] = useState(false);
    const defaultConfig = {
        message: "Bạn nhận được 1 cái nịt",
        onConfirm: () => 1,
        buttonText: t("Agree")
    }

    const [config, setConfig] = useState({ ...defaultConfig });
    useEffect(() => {
        messageDialogRef.current = {
            show: (config) => {
                setConfig({
                    ...defaultConfig,
                    ...config
                })
                setShow(true);
            },
            hide: () => {
                setShow(false);
                setConfig({ ...defaultConfig })
            }
        }
    }, []);

    return (
        <Modal
            visible={show}
            transparent={true}
        >
            <TouchableWithoutFeedback>
                <View className={className.overlay}>
                    <View className={className.content}>
                        <Text className={className.title}>{t("Notify")}</Text>
                        <Text
                            className={className.message}
                            style={{ textAlignVertical: "center" }}>
                            {config.message}
                        </Text>
                        <TouchableOpacity
                            className={className.button}
                            activeOpacity={0.8}
                            onPress={() => {
                                config.onConfirm();
                                setShow(false);
                            }}>
                            <Text className={className.buttonText}>
                                {config.buttonText.toUpperCase()}
                            </Text>
                        </TouchableOpacity>
                    </View>
                </View>
            </TouchableWithoutFeedback>
        </Modal>
    )
}

export default MessageDialog
export { messageDialogRef }