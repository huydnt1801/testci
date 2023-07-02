import { Pressable, Text, View } from "react-native";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import _ from "lodash"

import className from "./className";

const ButtonRow = ({ title, iconLeft, iconRight, classNames, onPress }) => {
    return (
        <Pressable
            className={className.wrapper + " " + classNames}
            onPress={onPress}>
            <View className={className.left}>
                {_.isObject(iconLeft) ? iconLeft : ""}
                <Text className={className.title}>{title}</Text>
            </View>
            {_.isObject(iconRight) ? iconRight : ""}
        </Pressable>
    )
}

export default ButtonRow