import { Text, View, Button, StyleSheet, TouchableOpacity, TextInput } from "react-native";
import MapView from "react-native-maps"
import className from "./className"
import { faCircleXmark, faLocationDot } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import MarkerSelect from "./components/MarkerSelect";


const SelectLocationOnMap = () => {

    const [isChanging, setIsChanging] = useState(false);

    const defaultPosition = {
        latitude: 21.0285,
        longitude: 105.8542
    };

    const { t } = useTranslation();
    const [searchInput, setSearchInput] = useState("");

    return (
        <View className={className.container}>
            <MapView
                style={StyleSheet.absoluteFillObject}
                initialRegion={{
                    latitude: defaultPosition.latitude,
                    longitude: defaultPosition.longitude,
                    latitudeDelta: 0.0522,
                    longitudeDelta: 0.0221
                }}
                showsBuildings={false}
                showsPointsOfInterest={false}
                showsUserLocation={false}
                onRegionChangeComplete={(region) => {
                    setIsChanging(false);
                    console.log(region);
                    console.log("End change");
                }}
                onRegionChange={() => {
                    setIsChanging(true);
                }}
            />
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
                        onPress={1}>
                        <FontAwesomeIcon
                            icon={faCircleXmark}
                            size={18}
                            style={{ color: "rgb(75 85 99)" }} />
                    </TouchableOpacity>
                )}
            </View>
            <MarkerSelect
                isChanging={isChanging} />
        </View>
    )
}

export default SelectLocationOnMap