<?php

namespace App\Http\Controllers;

use App\Services\ContactFormResponseService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class ContactFormResponseApiController
{
    private $_service;

    public function __construct(ContactFormResponseService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request)
    {
        $search = $request->query('search', '');
        $perPage = $request->query('per_page', 10);
        $sortBy = $request->query('sort_by');
        $sortDirection = $request->query('desc', 'desc');

        return response()->json($this->_service->all(perPage:$perPage, search: $search, sortBy: $sortBy, sortDirection: $sortDirection), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $sanitizedEmail = filter_var($request->input('email'), FILTER_SANITIZE_EMAIL);
        $sanitizedSubject = strip_tags($request->input('subject'));
        $sanitizedMessage = strip_tags($request->input('message'));

        $contactFormResponse = $this->_service->add([
            'email' => $sanitizedEmail,
            'subject' => $sanitizedSubject,
            'message' => $sanitizedMessage
        ]);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $contactFormResponse];
    }
}
