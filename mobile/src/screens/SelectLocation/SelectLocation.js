import { useTranslation } from "react-i18next";
import {
    Pressable,
    Text,
    TouchableOpacity,
    View
} from "react-native";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faArrowLeftLong, faCircleXmark, faLocationDot } from "@fortawesome/free-solid-svg-icons";

import className from "./className";
import { useNavigation } from "@react-navigation/native";
import { TextInput } from "react-native";
import { useState, useEffect } from "react";
import { useDebounce } from "../../hooks";
import axios from "axios";
import queryString from "query-string"
import Header from "../../components/Header";

const data_ = [
    {
        id: 1,
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
    {
        id: 2,
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
    {
        id: 3,
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội 12312  2123 2",
        distance: "5.9KM"
    },
    {
        id: 4,
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội sadasdsadsad",
        distance: "5.9KM"
    },
    {
        id: 5,
        text: '121 Lê Lợi',
        value: "121 Lê Lợi, Nguyễn Trãi, Hà Đông, Hà Nội",
        distance: "5.9KM"
    },
]

const API_KEY = 'AIzaSyBuExS3x9vzjyXhmDFxc8ZNWtcXHf5Kulg'


const SelectLocation = () => {

    const { t } = useTranslation();
    const navigation = useNavigation();

    const [data, setData] = useState([
        ...data_
    ])

    const [searchInput, setSearchInput] = useState("");
    const debounce = useDebounce(searchInput, 500);

    const getData = async () => {
        const result = await fetch(`https://maps.googleapis.com/maps/api/place/autocomplete/json?key=${API_KEY}&components=country:VN&input=${queryString.stringify(debounce)}`).then(i => i.json());
        console.log(JSON.stringify(result));
    }

    useEffect(() => {
        // getData()
    }, [debounce])


    return (
        <View className={className.container}>
            <Header
                title={t("Destination")}
                onPressBack={() => navigation.goBack()}
                componentRight={
                    <TouchableOpacity
                        activeOpacity={0.5}>
                        <Text className={className.map}>{t("SelectFromMap")}</Text>
                    </TouchableOpacity>
                } />
            <View className={className.search}>
                <FontAwesomeIcon
                    icon={faLocationDot}
                    size={24}
                    style={{ color: "rgb(234 179 8)", marginRight: 8 }} />
                <TextInput
                    className={className.input}
                    value={searchInput}
                    placeholder={t("EnterDestination")}
                    onChangeText={text => setSearchInput(text)} />
                {searchInput.length > 0 && (
                    <TouchableOpacity
                        style={{ padding: 4 }}
                        activeOpacity={0.5}
                        onPress={() => setSearchInput("")}>
                        <FontAwesomeIcon
                            icon={faCircleXmark}
                            size={18}
                            style={{ color: "rgb(75 85 99)" }} />
                    </TouchableOpacity>
                )}
            </View>
            <View className={className.placeSuggest}>
                {data.map((item, index) => (
                    <Pressable
                        key={item.id ?? index}
                        className={className.item}>
                        <View className={className.iconBorder}>
                            <FontAwesomeIcon
                                icon={faLocationDot} />
                        </View>
                        <View className={className.center}>
                            <Text className={className.textBold}>
                                {item.text}
                            </Text>
                            <Text
                                className={className.textLight}
                                numberOfLines={1}
                                lineBreakMode={"tail"}>
                                {item.value}
                            </Text>
                        </View>
                        <Text className={className.distance}>
                            {item.distance}
                        </Text>
                    </Pressable>
                ))}
            </View>
        </View>
    );
}

export default SelectLocation