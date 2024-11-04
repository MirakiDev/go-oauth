export const getMe = async () =>{ 
    const response = await fetch("http://localhost:8080/auth/me", {
        credentials: "include",
    });
    const data = await response.json();
    return data.RawData
}


