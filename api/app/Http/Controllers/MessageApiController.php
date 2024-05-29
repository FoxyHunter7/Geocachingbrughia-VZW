<?php

namespace App\Http\Controllers;

use App\Services\MessageService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class MessageApiController
{
    private $_service;

    public function __construct(MessageService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request, $isAdmin)
    {
        $language = $request->get('lang');
        $search = $request->query('search', '');
        $perPage = $request->query('per_page', 10);
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return $this->_service->all($language, $perPage, $search, $sortBy, $sortDirection, $isAdmin);
    }

    public function getByIdAdmin(Request $request, int $id)
    {
        $language = $request->get('lang');
        return response()->json($this->_service->getById($id, $language, true), Response::HTTP_OK);
    }

    public function allPublic (Request $request)
    {
        return response()->json($this->all($request, false), Response::HTTP_OK);
    }

    public function allAdmin (Request $request)
    {
        return response()->json($this->all($request, true), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $message = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $message];
    }

    public function update(Request $request, int $id)
    {
        $data = $request->all();
        $message = $this->_service->update($id, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $message];
    }

    public function delete(int $id)
    {
        $deletedMessage = $this->_service->delete($id);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedMessage];
    }
}
