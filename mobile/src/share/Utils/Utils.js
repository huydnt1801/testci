import CheckNetwork, { checkNetworkRef } from "./components/CheckNetwork";
import Loading, { loadingRef } from "./components/Loading";
import MessageDialog, { messageDialogRef } from "./components/MessageDialog";
import Toast, { toastRef } from "./components/Toast";

const wait = (ms) => new Promise(e => setTimeout(e, ms));

const Utils = {
    data: {},
    /**
     * Show a toast to screen
     * @param {String} message: message of toast
     * @param {Number | undefined} duration: time exist ms
     */
    toast: (message = "", duration = 2500) => {
        toastRef.current.hide();
        toastRef.current.show(message, duration);
    },
    /**
     * Show message dialog
     * @typedef config
     * @property {String | undefined} message
     * @property {() => void | undefined} onConfirm
     * @property {String | undefined} buttonText
     * @param {config} config 
     */
    showMessageDialog: (config) => {
        messageDialogRef.current.hide();
        messageDialogRef.current.show(config)
    },
    hideMessageDialog: () => {
        messageDialogRef.current.hide();
    },
    showCheckNetwork: () => {
        checkNetworkRef.current.hide();
        checkNetworkRef.current.show();
    },
    hideCheckNetwork: () => {
        checkNetworkRef.current.hide();
    },
    showLoading: () => {
        loadingRef.current.hide();
        loadingRef.current.show();
    },
    hideLoading: () => {
        loadingRef.current.hide();
    },
    /**
     * wait ms before next action
     * @param {Number} ms 
     */
    wait: async (ms) => {
        await wait(ms);
    },
}

const UtilComponents = () => {

    return (
        <>
            <Toast />
            <MessageDialog />
            <CheckNetwork />
            <Loading />
        </>
    );
}

export default Utils
export { UtilComponents }