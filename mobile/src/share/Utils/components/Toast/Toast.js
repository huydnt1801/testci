import { createRef, useEffect, useState } from "react";
import { StyleSheet, Text, View } from "react-native";

const toastRef = createRef();

const Toast = () => {

    const [show, setShow] = useState(false);
    const [message, setMessage] = useState("");

    useEffect(() => {
        toastRef.current = {
            show(message = "", duration = options.time.LONG) {
                setMessage(message);
                setTimeout(() => setShow(true));
                setTimeout(() => setShow(false), duration);
            },
            hide() {
                setShow(false);
            }
        }
    }, []);

    return (
        <View style={styles.wrapper}>
            {show &&
                <View style={styles.content}>
                    <Text style={styles.message}>{message}</Text>
                </View>
            }
        </View>
    );
}

const styles = StyleSheet.create({
    wrapper: {
        position: "absolute",
        bottom: "10%",
        width: "100%",
        zIndex: 1,
        alignItems: "center"
    },
    content: {
        backgroundColor: "rgba(0, 0, 0, 0.8)",
        padding: 10,
        borderRadius: 8,
        maxWidth: "80%"
    },
    message: {
        color: "white"
    }
})

const options = {
    time: {
        LONG: 4000,
        SHORT: 2500
    },
    location: {
        TOP: 1,
        CENTER: 2,
        BOTTOM: 3
    }
}

export default Toast
export { toastRef }