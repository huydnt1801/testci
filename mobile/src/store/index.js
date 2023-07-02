import { configureStore, applyMiddleware, compose } from "@reduxjs/toolkit";

import accountReducer from "../slices/Account";

const reducer = {
    account: accountReducer
}

const logger = (_store) => (next) => (action) => {
    return next(action);
}

const store = configureStore({
    reducer: reducer,
    devTools: true,
    enhancers: [compose(applyMiddleware(logger))],
    middleware: (getDefaultMiddleware) => getDefaultMiddleware({
        serializableCheck: false
    })
});

export default store