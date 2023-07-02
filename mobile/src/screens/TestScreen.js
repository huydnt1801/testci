import { useNavigation } from "@react-navigation/native";
import { Button, View, Text, ScrollView } from "react-native"
import { StatusBar } from "react-native";


const TestScreen = () => {

    const navigation = useNavigation();

    return (
        <View className="flex-1 flex-col pt-4">
            <StatusBar />
            <Text className="text-center mb-4 font-semibold">TestScreen</Text>
            <ScrollView>
                <Button
                    title="Trang chủ (Home)"
                    onPress={() => navigation.navigate("Home")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Đăng nhập (Login)"
                    onPress={() => navigation.navigate("Login")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Nhập MK hoặc OPT (PasswordOTP)"
                    onPress={() => navigation.navigate("PasswordOTP")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Cài đặt chuyến đi (TripSetting)"
                    onPress={() => navigation.navigate("TripSetting")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Cài đặt (Setting)"
                    onPress={() => navigation.navigate("Setting")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Thanh toán (Payment)"
                    onPress={() => navigation.navigate("Payment")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Tìm kiếm địa chỉ (SelectlLocation)"
                    onPress={() => navigation.navigate("SelectLocation")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Tìm kiếm trên bản đồ (SelectLocationOnMap)"
                    onPress={() => navigation.navigate("SelectLocationOnMap")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Đường đi (TripDirection)"
                    onPress={() => navigation.navigate("TripDirection")} />
                <View style={{ height: 12 }}></View>
                <Button
                    title="Ngôn ngữ (Language)"
                    onPress={() => navigation.navigate("Language")} />
                <View style={{ height: 12 }}></View>
            </ScrollView>
        </View>
    );
}
export default TestScreen