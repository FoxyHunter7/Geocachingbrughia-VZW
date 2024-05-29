<?php

namespace App\Http\Middleware;

use App\Services\LanguageService;
use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;
use Illuminate\Support\Facades\App;

class HandleLanguage
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $locale = App::getLocale();
        $language = $request->input("lang", $locale);
        $languageService = App(LanguageService::class);

        if ($languageService->checkIfExists($language))
        {
            if ($language != $locale)
            App::setLocale($locale);

            $request->attributes->set('lang', $language);
            return $next($request);
        }

        return response()->json(['error' => 'Language code: "' . $language . '" is not supported'], Response::HTTP_NOT_FOUND);
    }
}
