import { NavigationContainer } from "@react-navigation/native";
import { createNativeStackNavigator } from "@react-navigation/native-stack";
import { } from "nativewind";
import { I18nextProvider } from "react-i18next";
import { Provider } from "react-redux";

import Splash from "./src/screens/Splash";
import Home from "./src/screens/Home";
import Login from "./src/screens/Login";
import PasswordOTP from "./src/screens/PasswordOTP";
import Payment from "./src/screens/Payment";
import Setting from "./src/screens/Setting";
import TripSetting from "./src/screens/TripSetting";
import Mapp from "./src/screens/Mapp";
import SelectLocation from "./src/screens/SelectLocation";
import SelectLocationOnMap from "./src/screens/SelectLocationOnMap";
import TripDirection from "./src/screens/TripDirection";
import store from "./src/store";
import i18next from "./src/share/Language";
import { UtilComponents } from "./src/share/Utils";
import TestScreen from "./src/screens/TestScreen";
import Language from "./src/screens/Language";
import ChangePassword from "./src/screens/ChangePassword";

const Stack = createNativeStackNavigator();

const screens = [
    { name: "Splash", component: Splash },
    { name: "TestScreen", component: TestScreen },
    { name: "Home", component: Home },
    { name: "Login", component: Login },
    { name: "Payment", component: Payment },
    { name: "Setting", component: Setting },
    { name: "TripSetting", component: TripSetting },
    { name: "PasswordOTP", component: PasswordOTP },
    { name: "SelectLocation", component: SelectLocation },
    { name: "SelectLocationOnMap", component: SelectLocationOnMap },
    { name: "TripDirection", component: TripDirection },
    { name: "Mapp", component: Mapp },
    { name: "Language", component: Language},
    { name: "ChangePassword", component: ChangePassword },

]

const App = () => {
    return (
        <Provider store={store}>
            <I18nextProvider i18n={i18next}>
                <NavigationContainer>
                    <Stack.Navigator>
                        {screens.map(item => (
                            <Stack.Screen
                                key={item.name}
                                name={item.name}
                                component={item.component}
                                options={{ headerShown: false }} />
                        ))}
                    </Stack.Navigator>
                </NavigationContainer>
                <UtilComponents />
            </I18nextProvider>
        </Provider>
    );
}

export default App