import { useNavigation } from "@react-navigation/native";
import { View } from "react-native";
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs";
import { FontAwesomeIcon } from "@fortawesome/react-native-fontawesome";
import { faClockRotateLeft, faGift, faHome, faUser } from "@fortawesome/free-solid-svg-icons";

import HomePage from "./pages/HomePage";
import Account from "./pages/Account";
import History from "./pages/History";
import Reward from "./pages/Reward";
import { useSelector } from "react-redux";

const Tab = createBottomTabNavigator();

const pages = [
    {
        id: 0,
        name: "HomePage",
        component: HomePage,
        icon: faHome
    },
    {
        id: 1,
        name: "History",
        component: History,
        icon: faClockRotateLeft
    },
    {
        id: 2,
        name: "Reward",
        component: Reward,
        icon: faGift
    },
    {
        id: 3,
        name: "Account",
        component: Account,
        icon: faUser
    },
]

const Home = () => {

    const navigation = useNavigation();

    const { account } = useSelector(state => state.account);
    console.log(account);

    return (
        <Tab.Navigator
            screenOptions={{
                tabBarShowLabel: false,
                tabBarStyle: { height: 60 },
                tabBarHideOnKeyboard: true
            }}>
            {pages.map(item => (
                <Tab.Screen
                    key={item.id}
                    name={item.name}
                    component={item.component}
                    options={{
                        headerShown: false,
                        tabBarIcon: ({ focused }) => (
                            <View>
                                <FontAwesomeIcon
                                    className="bg-gray-300"
                                    icon={item.icon}
                                    size={24}
                                    style={{
                                        color: focused ? "rgb(234 179 8)" : " rgb(209 213 219)"
                                    }}
                                />
                            </View>
                        )
                    }}
                />
            ))}
        </Tab.Navigator>
    );
}

export default Home