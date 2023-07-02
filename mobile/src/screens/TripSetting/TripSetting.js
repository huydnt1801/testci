import { Text, View } from "react-native";
import { StackActions, useNavigation } from "@react-navigation/native";
import {
    ScrollView,
} from "react-native";
import Header from "../../components/Header";
import { useTranslation } from "react-i18next";
import className from "./className";
import ButtonRow from "./ButtonRow";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faFileInvoice,faLocation, faMessage, faPhone,} from "@fortawesome/free-solid-svg-icons";
import Utils from "../../share/Utils";

const TripSetting = () => {
    const navigation = useNavigation();
    const { t } = useTranslation();
    return (
        <View style={{ flex: 1 }}>
            <View >
                <Header
                    onPressBack={() => navigation.goBack()}
                    title={t("TripSetting")} />
                    
            </View>
            <ScrollView>

                <ButtonRow
                    classNames={`mt-2`}
                    title={t("SavedLocations")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <FontAwesomeIcon
                            icon={faLocation}
                            size={18}
                            style={{
                                color: "rgb(107 114 128)",
                                marginRight: 8
                            }} />}

                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("PresentMessages")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <FontAwesomeIcon
                            icon={faMessage}
                            size={18}
                            style={{
                                color: "rgb(107 114 128)",
                                marginRight: 8
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("AddEmergencyContacts")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <FontAwesomeIcon
                            icon={faPhone}
                            size={18}
                            style={{
                                color: "rgb(107 114 128)",
                                marginRight: 8
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("InvoiceInformation")}
                    iconLeft={
                        <FontAwesomeIcon
                            icon={faFileInvoice}
                            size={18}
                            style={{
                                color: "rgb(107 114 128)",
                                marginRight: 8
                            }} />}
                    onPress={() => Utils.toast("Coming Soon!")}
                />

            </ScrollView>
        </View>
    );
}

export default TripSetting