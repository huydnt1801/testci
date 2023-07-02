import { useState, useRef } from "react";
import { Text, View, Pressable, TextInput, StatusBar } from "react-native"
import { useTranslation } from "react-i18next";
import { StackActions, useNavigation, useRoute } from "@react-navigation/native";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useDispatch } from "react-redux";

import className from "./className";
import Header from "../../components/Header";
import Utils from "../../share/Utils";
import Api from "../../api";
import { setAccount } from "../../slices/Account";

const types = {
    ONLY_PASSWORD: 1,
    ONLY_OPT: 2,
    PASSWORD_AND_NAME: 3
}


/**
 * @typedef prop
 * @property {String | undefined} value
 * @property {() => void | undefined} onChangText
 * @param {prop} param
 * @returns 
 */
const InputRow = ({ value, onChangText, ref }) => {

    const _value = value ?? "";
    const renderList = ["", "", "", "", "", ""];
    for (let i = 0; i < _value.length; i++) {
        renderList[i] = _value[i];
    }

    const _ref = ref ?? useRef(null);

    return (
        <View className={className.inputWrapper}>
            <TextInput
                ref={_ref}
                style={{ position: 'absolute', width: "20%", zIndex: -1, opacity: 0, height: 40 }}
                cursorColor={"rgba(0,0,0,0)"}
                value={value}
                maxLength={6}
                onChangeText={onChangText}
                keyboardType={"number-pad"} />
            {renderList.map((item, index) => (
                <Pressable
                    key={index}
                    className={className.block + (index == value.length && " border-blue-400")}
                    onPress={() => {
                        _ref.current.blur()
                        _ref.current.focus();
                    }}>
                    <Text
                        className={className.inputItem}>
                        {item}
                    </Text>
                </Pressable>
            ))}
        </View>
    )
}

const PasswordOTP = () => {

    const { type: type_, userPhone } = useRoute().params;
    const [type, setType] = useState(type_ ?? types.PASSWORD_AND_NAME);

    const { t } = useTranslation();
    const navigation = useNavigation();

    const dispatch = useDispatch();

    const [passwordOrOTP, setPasswordOrOTP] = useState("");
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const titles = {
        [types.ONLY_PASSWORD]: t("EnterPassword"),
        [types.ONLY_OPT]: t("EnterOTP"),
        [types.PASSWORD_AND_NAME]: t("Register"),
    }
    const messages = {
        [types.ONLY_PASSWORD]: t("YouAreLoginWithPhone"),
        [types.ONLY_OPT]: t("EnterOTPSendToPhone"),
        [types.PASSWORD_AND_NAME]: t("EnterYourInformation"),
    }

    const isButtonActive = (() => {
        if (type == types.ONLY_PASSWORD) {
            if (passwordOrOTP.length == 6) {
                return true;
            }
            return false;
        }
        if (type == types.ONLY_OPT) {
            if (passwordOrOTP.length == 6) {
                return true;
            }
            return false;
        }
        if (type == types.PASSWORD_AND_NAME) {
            if (username.trim().length >= 6 && password.length == 6) {
                return true;
            }
            return false;
        }
    })();

    const handleClick = async () => {
        if (type == types.PASSWORD_AND_NAME) {
            Utils.showLoading();
            await Utils.wait(300);
            const result = await Api.account.register(userPhone, username, password);
            Utils.hideLoading();
            if (result.data) {
                Utils.data["token"] = null;
                Utils.data["token"] = result.data.data;
                setType(types.ONLY_OPT);
            }
        }
        if (type == types.ONLY_OPT) {
            Utils.showLoading();
            const result = await Api.account.confirmOTP(Utils.data.token, passwordOrOTP);
            await Utils.wait(300);
            if (result.result == Api.ResultCode.SUCCESS) {
                const login = await Api.account.login(userPhone, password);
                await Utils.wait(300);
                if (login.result == Api.ResultCode.SUCCESS) {
                    dispatch(setAccount(login.data.data));
                    try {
                        await AsyncStorage.setItem("account", JSON.stringify(login.data.data))
                    } catch (error) { }
                    Utils.hideLoading();
                    navigation.dispatch(StackActions.replace("Home"))
                }
            }
            else {
                Utils.hideLoading();
                Utils.showMessageDialog({
                    message: t("WrongOTP")
                });
            }
        }
        if (type == types.ONLY_PASSWORD) {
            Utils.showLoading();
            const login = await Api.account.login(userPhone, passwordOrOTP);
            await Utils.wait(300);
            if (login.result == Api.ResultCode.SUCCESS) {
                dispatch(setAccount(login.data.data));
                try {
                    await AsyncStorage.setItem("account", JSON.stringify(login.data.data))
                } catch (error) { }
                Utils.hideLoading();
                navigation.dispatch(StackActions.replace("Home"))
            }
            else {
                Utils.hideLoading();
                Utils.showMessageDialog({
                    message: t("WrongPassword")
                });
            }
            Utils.hideLoading();
        }
    }

    return (
        <View className={className.container}>
            <StatusBar />
            <Header
                border={false}
                onPressBack={() => navigation.goBack()} />
            <View className={className.content}>
                {(type == types.ONLY_OPT || type == types.ONLY_PASSWORD) && (
                    <>
                        <Text className={className.title}>{titles[type]}</Text>
                        <Text>{messages[type]}</Text>
                        <Text className={className.phone}>{`+84 ${userPhone?.slice(1)}`}</Text>
                        <InputRow
                            value={passwordOrOTP}
                            onChangText={e => {
                                setPasswordOrOTP(
                                    e.replace(",", "")
                                        .replace(" ", "")
                                        .replace(".", "")
                                        .replace("-", "")
                                )
                            }}
                        />
                    </>
                )}
                {type == types.PASSWORD_AND_NAME && (
                    <>
                        <Text className={className.title}>{titles[type]}</Text>
                        <Text>{messages[type]}:</Text>
                        <Text className={className.enterName}>{t("EnterYourName")}:</Text>
                        <View className={className.border}>
                            <TextInput
                                value={username}
                                // placeholderTextColor={}
                                className={className.username + (username.length > 0 && " text-lg font-bold")}
                                onChangeText={text => setUsername(text)}
                            // placeholder={t("EnterYourName")} 
                            />
                        </View>
                        <Text className={className.enterName}>{t("EnterPassword")}:</Text>
                        <InputRow
                            value={password}
                            onChangText={e => setPassword(e)}
                        />
                    </>
                )}
            </View>
            <View className={className.bottom}>
                <Pressable
                    className={isButtonActive ? className.activate : className.deactivate}
                    disabled={!isButtonActive}
                    onPress={handleClick}>
                    <Text className={className.buttonText}>{t("Continue")}</Text>
                </Pressable>
            </View>
        </View>
    );
}

export default PasswordOTP
export { types }