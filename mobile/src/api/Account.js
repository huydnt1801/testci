import baseApi from "./baseApi";

const login = (phone, password) => {
    const url = "/accounts/login";
    const param = {
        phoneNumber: phone,
        password
    }
    return baseApi.post(url, param);
}
const confirmOTP = (token, otp) => {
    const url = "/accounts/register/confirm";
    const param = {
        token,
        type: "submit-otp",
        otp
    }
    return baseApi.post(url, param);
}

/**
 * 
 * @param {String} phone
 * @param {String} username 
 * @param {String} password 
 * @returns {Promise<void>}
 */
const register = (phone, username, password) => {
    const url = "/accounts/register";
    const param = {
        phoneNumber: phone,
        fullName: username,
        password: password,
        password2: password
    }
    return baseApi.post(url, param);
}

/**
 * 
 * @param {String} phone 
 * @returns {Promise<void>}
 */
const checkPhone = (phone) => {
    const url = "/accounts/phone";
    const param = {
        phoneNumber: phone
    }
    return baseApi.post(url, param);
}

const accountApi = {
    login,
    checkPhone,
    register,
    confirmOTP
}

export default accountApi