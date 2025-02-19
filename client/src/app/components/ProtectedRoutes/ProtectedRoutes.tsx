import React, { useEffect } from 'react';
import { Outlet, Navigate } from 'react-router-dom';
import Cookies from 'js-cookie';

export const ProtectedRoutes = () => {
    const auth = {'token': false};
    if (Cookies.get('khub-login-session-store')) {
        auth.token = true;
    }
    return(
        auth.token ? <Outlet/> : <Navigate to={'/sign-in'}/>
    );
};

export const LoginFlow = () => {
    useEffect(() => {
        const loginUrl = 
            process.env.REACT_APP_API_URL !== undefined && 
            process.env.REACT_APP_API_URL !== '' && 
            process.env.REACT_APP_API_URL !== null ? `${process.env.REACT_APP_API_URL}/login` : `${window.location.origin.toString()}/login`;
        window.location.replace(loginUrl);
    }, []);

    return <div>authenticating. . . </div>;
};