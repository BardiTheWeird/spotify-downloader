import React from "react";

export async function BackendPull() {
    const [isResourses, updateIsResourses] = React.useState();
    updateIsResourses(await fetch());

    return <></>
}