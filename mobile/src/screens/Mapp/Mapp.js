import { Text, View, Button, StyleSheet } from "react-native";
import MapView from "react-native-maps"

const Mapp = () => {

    return (
        <View style={StyleSheet.absoluteFillObject}>
            <MapView style={StyleSheet.absoluteFillObject} />
            <View style={{ position: 'absolute', top: 100, left: 50 }} />
        </View>
    );
}

const style = StyleSheet.absoluteFill

export default Mapp