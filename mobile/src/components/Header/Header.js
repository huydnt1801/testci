import { View, Text, TouchableOpacity, Image, Pressable } from "react-native";

import className from "./className";
import { iconArrowBack } from "../Icon";

const Header = ({ title, onPressBack, componentRight, border = true }) => {

    return (
        <View className={className.wrapper + (border ? " border-b border-gray-300" : "")}>
            <View className={className.left}>
                <TouchableOpacity
                    activeOpacity={0.5}
                    onPress={onPressBack}>
                    <Image
                        source={iconArrowBack}
                        style={{
                            width: 40,
                            height: 40,
                            tintColor: "black"
                        }} />
                </TouchableOpacity>
            </View>
            <View className={className.center}>
                <Text className={className.title}>{title}</Text>
            </View>
            <View className={className.right}>
                {componentRight}
            </View>
        </View>
    )
}

export default Header