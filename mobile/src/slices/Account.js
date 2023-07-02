import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    account: null
}

const accountSlice = createSlice({
    name: "account",
    initialState: initialState,
    reducers: {
        setAccount(state, { payload }) {
            state.account = payload
        },
    },
    extraReducers: (builder) => { }
});

const { reducer: accountReducer, actions } = accountSlice;
const { setAccount } = actions;

export default accountReducer
export { setAccount }