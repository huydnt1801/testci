import { View } from "react-native"
import className from "./className";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import { animated, useSpring } from "@react-spring/native";
import { memo } from "react";


const MarkerSelect = ({ isChanging }) => {

    const [style, spring] = useSpring(() => ({
        bottom: 20
    }));
    spring.start({
        bottom: isChanging ? 28 : 20,
        config: {
            duration: 100
        }
    });

    return (
        <View className={className.wrapper}>
            <animated.View
                className={className.top}
                style={style}>
                <View className={className.circle}>
                    <FontAwesomeIcon
                        icon={faHeart}
                        size={12}
                        style={{ color: "white" }} />
                </View>
                <View className={className.line}></View>
            </animated.View>
            <View className={className.shadow} ></View>
        </View>
    );
}

export default memo(MarkerSelect)