import { ImageBackground, Pressable, TextInput, View, StyleSheet } from "react-native"
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";

import { banhmique } from "../../../../../../components/Image";
import className from "./className"
import { faLocationDot } from "@fortawesome/free-solid-svg-icons";
import bebanner from "../../../../../../assets/images/bebanner.png";

const Banner = ({ onPress }) => {

    return (
        <ImageBackground
            source={bebanner}
            // style={StyleSheet.absoluteFillObject}
            className={className.wrapper}>
            <Pressable
                className={className.border}>
                <FontAwesomeIcon
                    icon={faLocationDot}
                    size={18}
                    style={{ color: "rgb(234 179 8)", marginRight: 8 }} />
                <TextInput
                    placeholder={"Bạn muốn đi đâu?"} />
                <Pressable
                    className={className.overlay}
                    onPress={onPress}>
                </Pressable>
            </Pressable>
            <View className={className.devide}></View>
        </ImageBackground>
    );
}

export default Banner