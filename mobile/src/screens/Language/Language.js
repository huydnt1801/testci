import { Button } from "react-native";
import { Text, View } from "react-native";
import { useDispatch } from "react-redux";
import { setAccount } from "../../slices/Account";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { StackActions, useNavigation } from "@react-navigation/native";
import {
    ScrollView,
} from "react-native";
import Header from "../../components/Header";
import { useTranslation } from "react-i18next";
import className from "./className";
import ButtonRow from "./ButtonRow";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faCheck } from "@fortawesome/free-solid-svg-icons";
import Utils from "../../share/Utils";

const Language = () => {

    const dispatch = useDispatch();
    const navigation = useNavigation();
    const { i18n } = useTranslation();
    const language = i18n.language;
    const { t } = useTranslation();
    return (
        <View className={className.container}>
            <View >
                <Header
                    title={t("SelectLanguage")}
                    onPressBack={()=> navigation.goBack()} />
            </View>
            <ScrollView>

                <ButtonRow
                    classNames={`mt-2`}
                    title={t("English")}
                    onPress={() => i18n.changeLanguage("en")}
                    iconRight={
                        <FontAwesomeIcon
                            icon={faCheck}
                            size={18}
                            style={{
                                color: language == "en" ? "rgb(234 179 8)" : "white",
                                marginRight: 8
                            }} />}

                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("Vietnamese")}
                    onPress={() => i18n.changeLanguage("vi")}
                    iconRight={
                        <FontAwesomeIcon
                            icon={faCheck}
                            size={18}
                            style={{
                                color: language == "vi" ? "rgb(234 179 8)" : "white",
                                marginRight: 8
                            }} />}
                />

            </ScrollView>
        </View>
    );
}


export default Language