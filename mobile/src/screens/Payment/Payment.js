import { Button } from "react-native";
import { Text, View, Image} from "react-native";
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
import { faCreditCard} from "@fortawesome/free-solid-svg-icons";
import Utils from "../../share/Utils";
import { iconCake, iconMomo, iconShoppeePay, iconViettelMoney, iconZaloPay } from "../../components/Icon";

const Payment = () => {

    const { t } = useTranslation();
    const navigation = useNavigation();

    return (
        <View style={{ flex: 1 }}>
            <Header
                title={t("Payment")}
                onPressBack={() => navigation.goBack()} />

            <ScrollView>
            <View>
                <Text style={{fontWeight: "bold", fontSize:15, margin:10}} > Thêm hình thức thanh toán mới </Text>
            </View>
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("CakeByVPBank")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <Image
                            source={iconCake}
                            style={{
                                width: 40,
                                height: 40,
                                marginRight: 15,
                                tintColor: "black"
                            }} />}

                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("CreditCardDebitCard")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <FontAwesomeIcon
                            icon={faCreditCard}
                            size={30}
                            style={{
                                color: "rgb(107 114 128)",
                                marginRight: 20,
                                marginLeft: 5
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("ZaloPay")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <Image
                            source={iconZaloPay}
                            style={{
                                width: 30,
                                height: 30,
                                marginRight: 20,
                                marginLeft: 5,
                                tintColor: "black"
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("MoMo")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <Image
                            source={iconMomo}
                            style={{
                                width: 50,
                                height: 50,
                                marginRight:10,
                                marginLeft:-5,
                                tintColor: "black"
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("ShopeePay")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <Image
                            source={iconShoppeePay}
                            style={{
                                width: 50,
                                height: 30,
                                marginRight: 5,
                                tintColor: "black"
                            }} />}
                />
                <ButtonRow
                    classNames={`mt-2`}
                    title={t("ViettelMoney")}
                    onPress={() => Utils.toast("Coming Soon!")}
                    iconLeft={
                        <Image
                            source={iconViettelMoney}
                            style={{
                                width: 40,
                                height: 40,
                                marginRight:15,
                                tintColor: "black"
                            }} />}
                />

            </ScrollView>
        </View>
    );
}

export default Payment