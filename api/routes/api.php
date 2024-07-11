<?php

use App\Http\Controllers\ContactFormResponseApiController;
use App\Http\Controllers\EventApiController;
use App\Http\Controllers\GeocacheApiController;
use App\Http\Controllers\ImageController;
use App\Http\Controllers\JwtAuthController;
use App\Http\Controllers\LanguageApiController;
use App\Http\Controllers\MessageApiController;
use App\Http\Controllers\SocialApiController;
use App\Http\Controllers\StaticSiteContentController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider and all of them will
| be assigned to the "api" middleware group. Make something great!
|
*/

Route::group([
    'middleware' => ['throttle:api']
], function() {
    Route::group([
        'middleware' => ['language']
    ], function() {
        Route::get('events', [EventApiController::class, 'allPublic']);
        Route::get('home_events', [EventApiController::class, 'homePageEvents']);
        Route::get('messages', [MessageApiController::class, 'allPublic']);
    });

    Route::get('geocaches', [GeocacheApiController::class, 'allPublic']);
    Route::get('languages', [LanguageApiController::class, 'allPublic']);
    Route::get('static',[StaticSiteContentController::class, 'allPublic']);
    Route::get('socials', [SocialApiController::class, 'allPublic']);

    Route::get('images/{p1}/{filename}', [ImageController::class, 'get'])->name('images');

    Route::post('contact/form/responses', [ContactFormResponseApiController::class, 'add']);

    // Only needed initially, users are not allowed to self register.
    Route::post('register', [JwtAuthController::class, 'register']);
    Route::post('login', [JwtAuthController::class, 'login']);
    Route::get('denied', [JwtAuthController::class, 'accessDenied'])->name('denied');
});

Route::group([
    'middleware' => ['auth.csrf.jwt', 'auth:api']
], function() {
    Route::controller(JwtAuthController::class)->group(function () {
        Route::get('profile', 'profile');
        Route::get('refresh', 'refreshToken');
        Route::get('logout', 'logout');
        //Route::post('register', 'register'); To be enabled after making initial user
    });

    Route::prefix('admin')->group(function () {
        Route::controller(EventApiController::class)->group(function () {
            Route::get('events', 'allAdmin');
            Route::get('events/{id}', 'getByIdAdmin');
            Route::post('events', 'add');
            Route::put('events/{id}', 'update');
            Route::delete('events/{id}', 'delete');
        });

        Route::controller(GeocacheApiController::class)->group(function () {
            Route::get('geocaches', 'allAdmin');
            Route::get('geocaches/{id}', 'getByIdAdmin');
            Route::post('geocaches', 'add');
            Route::put('geocaches/{id}', 'update');
            Route::delete('geocaches/{id}', 'delete');
        });

        Route::controller(MessageApiController::class)->group(function () {
            Route::get('messages', 'allAdmin');
            Route::get('messages/{id}', 'getByIdAdmin');
            Route::post('messages', 'add');
            Route::put('messages/{id}', 'update');
            Route::delete('messages/{id}', 'delete');
        });

        Route::controller(LanguageApiController::class)->group(function () {
            Route::get('languages', 'allAdmin');
            Route::post('languages', 'add');
            Route::put('languages/{lang}', 'update');
            Route::delete('languages/{lang}', 'delete');
        });

        Route::controller(StaticSiteContentController::class)->group(function () {
            Route::get('static', 'allAdmin');
            Route::post('static', 'add');
            Route::put('static/{property}', 'update');
            Route::delete('static/{property}', 'delete');
        });

        Route::controller(SocialApiController::class)->group(function () {
            Route::get('socials', 'allAdmin');
            Route::post('socials', 'add');
            Route::put('socials/{id}', 'update');
            Route::delete('socials/{id}', 'delete');
        });

        Route::controller(ContactFormResponseApiController::class)->group(function () {
            Route::get('contact/form/responses', 'all');
        });
    });
});
