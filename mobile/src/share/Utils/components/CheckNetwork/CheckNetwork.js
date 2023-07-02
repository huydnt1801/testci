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
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";

const checkNetworkRef = createRef();

const CheckNetwork = () => {

    const { t } = useTranslation();
    const [show, setShow] = useState(false);

    useEffect(() => {
        checkNetworkRef.current = {
            show: () => {
                setShow(true);
            },
            hide: () => {
                setShow(false);
            }
        }
    }, []);


    return (
        <Modal
            visible={show}
            transparent={true}>
            <TouchableWithoutFeedback>
                <View className={className.overlay}>
                    <View className={className.content}>
                        <Text className={className.title}>{t("Notify")}</Text>
                        <FontAwesomeIcon
                            icon={faTriangleExclamation}
                            size={50}
                            style={{ color: "red", alignSelf: "center", marginTop: 12 }} />
                        <Text className={className.message}>{t("CanNotConnectInternet")}</Text>
                        <Text className={className.message}>{t("PleaseTryAgain")}</Text>
                        <TouchableOpacity
                            className={className.button}
                            activeOpacity={0.8}
                            onPress={() => {
                                setShow(false);
                            }}>
                            <Text className={className.buttonText}>
                                {t("Agree")}
                            </Text>
                        </TouchableOpacity>
                    </View>
                </View>
            </TouchableWithoutFeedback>
        </Modal>
    );
};

export default CheckNetwork
export { checkNetworkRef }