<?php

namespace App\Services;

use App\Models\User;
use Illuminate\Support\MessageBag;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Hash;
use Tymon\JWTAuth\Facades\JWTAuth;
use Nette\Utils\Random;

class JwtAuthService
{
    protected $_model;
    protected $_rules = [
        'register' =>  [
            'name' => 'required',
            'email' => 'required|email|unique:users',
            'password' => 'required|confirmed'
        ],
        'login' => [
            'email' => 'required|email',
            'password' => 'required'
        ]
    ];

    protected $_errors;

    public function __construct(User $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    protected function validate($data, $ruleset)
    {
        $this->_errors = new MessageBag();

        $validator = Validator::make($data, $this->_rules[$ruleset]);

        if ($validator->fails()) {
            $this->_errors = $validator->errors();
        }
    }

    public function hasErrors() {
        return $this->_errors->any();
    }

    public function getErrors() {
        return $this->_errors;
    }

    public function register($data)
    {
        $this->validate($data, "register");
        if ($this->hasErrors()) {
            return;
        }

        $data['password'] = Hash::make($data['password']);
        $this->_model->create($data);
    }

    public function login($data)
    {
        $this->validate($data, "login");
        if ($this->hasErrors()) {
            return;
        }

        $csrfLength = env("CSRF_TOKEN_LENGTH");
        $csrfToken = Random::generate($csrfLength);

        return [JWTAuth::claims(['X-XSRF-TOKEN' => $csrfToken])->attempt($data), $csrfToken];
    }

    public function refreshToken()
    {
        $csrfLength = env("CSRF_TOKEN_LENGTH");
        $csrfToken = Random::generate($csrfLength);

        return [JWTAuth::claims(['X-XSRF-TOKEN' => $csrfToken])->refresh(), $csrfToken];
    }
}
