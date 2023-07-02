import { Button, StyleSheet } from "react-native";
import { TextInput, View, text, setText } from "react-native";
import { useNavigation } from "@react-navigation/native";
import {
    ScrollView, 
} from "react-native";
import Header from "../../components/Header";
import { useTranslation } from "react-i18next";
import className from "./className";
import Utils from "../../share/Utils";

const ChangePassword = () => {

    const navigation = useNavigation();
    const { t } = useTranslation();
    return (
        <View className={className.container}>
            <View >
                <Header
                    title={t("ChangePassword")}
                    onPressBack={() => navigation.goBack()} />
            </View>
            <ScrollView>
                <TextInput
                    name="oldPassword"
                    placeholder={t("EnterOldPassword")}
                    keyboardType="numeric"
                    secureTextEntry={true}
                    style={{
                        backgroundColor: "white",
                        marginTop: 20,
                        height: 50,
                        paddingLeft: 20
                    }}
                />
                <TextInput
                    name="newPassword"
                    placeholder={t("EnterANewPassword")}
                    keyboardType="numeric"
                    secureTextEntry={true}
                    style={{
                        backgroundColor: "white",
                        marginTop: 20,
                        height: 50,
                        paddingLeft: 20
                    }}
                />
                <TextInput
                    name="confirmNewPassword"
                    placeholder={t("EnterYourNewPassword")}
                    keyboardType="numeric"
                    secureTextEntry={true}
                    style={{
                        backgroundColor: "white",
                        marginTop: 20,
                        height: 50,
                        paddingLeft: 20,
                        marginBottom:20
                    }}
                />
                <Button
                    title={t("ChangePassword")}
                    name="changePassword"
                    onPress={() => Utils.toast("Coming soon")}

                />

            </ScrollView>
        </View>
    );
}

export default ChangePassword