import {
    Pressable,
    Text,
    TouchableOpacity,
    View
} from "react-native";
import {
    ScrollView,
} from "react-native";
import { useTranslation } from "react-i18next";
import className from "./className";
import ButtonRow from "./ButtonRow";
import { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faBiking, faCableCar, faClock, faFingerprint, faKeyboard, faLanguage, faLocationDot, faPeopleArrows, faPeopleArrowsLeftRight, faPersonPraying, faPersonWalking } from "@fortawesome/free-solid-svg-icons";
import Utils from "../../../../share/Utils";

const data_ = [
    {
        id: 1,
        time: '11/2/2022 | 12:33',
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
    {
        id: 2,
        time: '12/2/2022 | 12:33',
        text: '122 Lê Lợi',
        value: "122 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
    {
        id: 3,
        time: '13/2/2022 | 12:33',
        text: '123 Lê Lợi',
        value: "123 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội ",
        distance: "5.9KM"
    },
    {
        id: 4,
        time: '14/2/2022 | 12:33',
        text: '124 Lê Lợi',
        value: "124 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội ",
        distance: "5.9KM"
    },
    {
        id: 5,
        time: '15/2/2022 | 12:33',
        text: '125 Lê Lợi',
        value: "125 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội, Nguyễn Trãi, Hà Đông, Hà Nội, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
]

const History = () => {

    const { t } = useTranslation();
    const [data, setData] = useState([
        ...data_
    ])
    return (
        <View className={className.container}>
            <View>
                <Text style={{ fontWeight: "bold", fontSize: 20, margin: 10, color: "rgb(234 179 8)" }} > {t("TripHistory")} </Text>
            </View>
            <ScrollView className={className.content}>
                {data.map((item, index) => (
                    <Pressable
                        key={item.id ?? index}
                        className={className.element}>
                        <ButtonRow
                            title={item.time}
                            classNames={`m-0`}
                            iconLeft={<FontAwesomeIcon
                                icon={faClock}
                                size={18}
                                style={{
                                    marginRight: 8,
                                    color:"rgb(234 179 8)"
                                }} />}
                        />
                        <ButtonRow
                            title={item.text}
                            classNames={`m-0`}
                            iconLeft={<FontAwesomeIcon
                                icon={faPersonWalking}
                                size={18}
                                style={{
                                    marginRight: 8,
                                    color: "rgb(234 179 8)"
                                }}
                            />}
                        />
                        <ButtonRow
                            title={item.value}
                            classNames={`m-0`}
                            iconLeft={<FontAwesomeIcon
                                icon={faLocationDot}
                                size={18}
                                style={{
                                    marginRight: 8,
                                    color: "rgb(234 179 8)"
                                }}
                            />}
                        />
                    </Pressable>
                ))}
            </ScrollView>
        </View>
    );
}

export default History