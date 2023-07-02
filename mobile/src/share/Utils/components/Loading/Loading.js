import { createRef, useEffect, useState } from "react";
import { View, ActivityIndicator, Modal } from "react-native";

const loadingRef = createRef();

const Loading = () => {

    const [show, setShow] = useState(false);

    useEffect(() => {
        loadingRef.current = {
            show: () => {
                setShow(true);
            },
            hide: () => {
                setShow(false);
            }
        }
    }, [])

    return (
        <Modal
            visible={show}
            transparent={true}
        // statusBarTranslucent={true}
        >
            <View className={`
                flex-1 bg-black/30 items-center justify-center px-6
            `}>
                <View className={'scale-[2]'}>
                    <ActivityIndicator
                        color={"rgb(234,179,8)"}
                    />
                </View>
            </View>
        </Modal>
    );
}

export default Loading
export { loadingRef }