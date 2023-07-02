import { useEffect, useState } from "react";
import { View, Image } from "react-native";
import { useSpring, animated } from "@react-spring/native";
import { useDispatch } from "react-redux";
import { useNavigation, StackActions } from "@react-navigation/native";

import { docker } from "../../components/Image";
import { setAccount } from "../../slices/Account";
import Utils from "../../share/Utils";
import AsyncStorage from "@react-native-async-storage/async-storage";

const Splash = () => {

    const navigation = useNavigation();
    const dispatch = useDispatch();

    const [style, spring] = useSpring(() => ({
        top: "100%",
        backgroundColor: "#d4ecf2"
    }));

    const rememberLogin = async () => {
        await Utils.wait(1000);
        Utils.showLoading();
        const account = await (async () => {
            try {
                const data = await AsyncStorage.getItem("account");
                return data;
            } catch (error) {
                return null;
            }
        })();
        await Utils.wait(300);
        Utils.hideLoading();
        if (account) {
            dispatch(setAccount(JSON.parse(account)))
            navigation.dispatch(StackActions.replace("Home"));
        }
        else {
            navigation.dispatch(StackActions.replace("Login"));
        }
    }

    useEffect(() => {
        spring.start({
            top: "0%",
            config: {
                duration: 700
            }
        });
        rememberLogin();
    }, []);

    return (
        <View className="flex-1" style={{ backgroundColor: "#d4ecf2" }}>
            <animated.View className="items-center justify-center absolute w-full h-full" style={style}>
                <Image
                    source={docker}
                    resizeMode={"contain"}
                    style={{ width: "60%", marginBottom: 50 }}
                />
            </animated.View>
        </View>
    );
}

export default Splash