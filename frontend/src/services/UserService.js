import React from "react";

export const UserContext = React.createContext();

export function UserProvider(props) {
    const [user, updateUser] = React.useState();
    return (
        <UserContext.Provider value={[user, updateUser]}>
            {props.children}
        </UserContext.Provider>
    )
}
