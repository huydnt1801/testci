import { createRef, useEffect, useState } from "react"
import { Modal } from "react-native";

const confirmDialogRef = createRef();


const ConfirmDialog = () => {
    const [show, setShow] = useState(false);

    useEffect(() => {
        confirmDialogRef.current = {
            show: () => {

            }
        }
    });

    return (
        <Modal visible={show}>

        </Modal>
    )
}

export default ConfirmDialog
export {
    confirmDialogRef
}