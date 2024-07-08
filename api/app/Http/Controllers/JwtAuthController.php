<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Services\JwtAuthService;
use Nette\Utils\Random;
use Tymon\JWTAuth\Facades\JWTAuth;

class JwtAuthController
{
    private $_service;

    public function __construct(JwtAuthService $service)
    {
        $this->_service = $service;
    }

    public function accessDenied()
    {
        return response()->json(["status" => false, "access_denied" => "You are not autherised to access this restricted endpoint"]);
    }

    public function register(Request $request)
    {
        $data = $request->all();
        $this->_service->register($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return response()->json([
            "status" => true,
            "message" => "User registered successfully"
        ]);
    }

    public function login(Request $request)
    {
        $data = $request->all();
        [$token, $csrfToken] = $this->_service->login($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        if(empty($token)){
            return response()->json([
                "status" => false,
                "message" => "Invalid details"
            ]);
        }

        $ttl = env("JWT_COOKIE_TTL");
        $tokenCookie = cookie("token", $token, $ttl);
        $csrfCookie = cookie("X-XSRF-TOKEN", $csrfToken, $ttl);

        return response()->json(["message" => "User logged in succcessfully"])
            ->withCookie($tokenCookie)
            ->withCookie($csrfCookie);
    }

    public function profile()
    {

        $userdata = auth()->user();

        return response()->json([
            "status" => true,
            "message" => "Profile data",
            "data" => $userdata
        ]);
    }

    public function refreshToken()
    {
        [$token, $csrfToken] = $this->_service->refreshToken();

        $ttl = env("JWT_COOKIE_TTL");
        $tokenCookie = cookie("token", $token, $ttl);
        $csrfCookie = cookie("X-XSRF-TOKEN", $csrfToken, $ttl);

        return response()->json(["message" => "Token refresh sucessful"])
            ->withCookie($tokenCookie)
            ->withCookie($csrfCookie);
    }

    public function logout()
    {

        auth()->logout();

        return response()->json([
            "status" => true,
            "message" => "User logged out successfully"
        ]);
    }
}
