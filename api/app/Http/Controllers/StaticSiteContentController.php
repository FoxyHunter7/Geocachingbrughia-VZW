<?php

namespace App\Http\Controllers;

use App\Services\StaticSiteContentService;
use Illuminate\Http\Request;

class StaticSiteContentController
{
    private $_service;

    public function __construct(StaticSiteContentService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request)
    {
        $search = $request->query('search', '');
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return response()->json($this->_service->all(search: $search, sortBy: $sortBy, sortDirection: $sortDirection));
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $event = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $event];
    }

    public function update(Request $request, int $id)
    {
        $data = $request->all();
        $event = $this->_service->update($id, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $event];
    }

    public function delete(int $id)
    {
        $deletedEvent = $this->_service->delete($id);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedEvent];
    }
}
