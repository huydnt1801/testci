import { Text, View, StatusBar, Pressable, Button, ScrollView, Image, StyleSheet } from "react-native";
import { useNavigation } from "@react-navigation/native"

import className from "./className";
import Banner from "./components/Banner/Banner";
import { bebanner, bevoucher, blitzcrank } from "../../../../components/Image";
import { useSelector } from "react-redux";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faCheck, faDriversLicense, faMotorcycle, faPerson } from "@fortawesome/free-solid-svg-icons";
import MapView from "react-native-maps"
const HomePage = () => {

    const name = "Trương Quang Phú";
    const navigation = useNavigation();
    const { account } = useSelector(state => state.account)
    return (
        <View className={className.container}>

            <StatusBar />
            <View className={className.header}>
                <Text className={className.welcome}>
                    {`Chào ${account?.full_name}!`}
                </Text>
                <Pressable>

                </Pressable>

            </View>
            <View>

                <Banner
                    onPress={() => navigation.navigate("SelectLocation")} />
                <View className={className.action}>
                    <Pressable className={className.button}>

                        <FontAwesomeIcon
                            icon={faMotorcycle}
                            size={40}
                            style={{
                                color: "rgb(234 179 8)",
                                marginRight: 8
                            }} />
                        {/* <Image
                            className={className.imageButton}
                            source={blitzcrank} /> */}
                        <Text className={className.textButton}>{"Tìm xe"}</Text>
                    </Pressable>
                    <Pressable className={className.button}>
                        <FontAwesomeIcon
                            icon={faPerson}
                            size={40}
                            style={{
                                color: "rgb(234 179 8)",
                                marginRight: 8
                            }} />
                        {/* <Image
                            className={className.imageButton}
                            source={blitzcrank} /> */}
                        <Text className={className.textButton}>{"Thuê tài xế"}</Text>
                    </Pressable>
                </View>
                <View>
                    <Image
                        style={{
                            resizeMode: 'repeat',
                            // height: 100,
                            // width: 200,
                            // display: 'flex'
                        }}
                        source={bevoucher} />
                </View>

            </View>
            {/* <Image
                style={{
                    resizeMode: 'center',
                    // height: 100,
                    // width: 200,
                    // display: 'flex'
                }}
                source={bevoucher} /> */}
        </View>
    );
}

export default HomePage